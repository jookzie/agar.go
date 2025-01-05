import React, { useEffect, useRef, useState, useCallback } from 'react';

interface Player {
	x: number;
	y: number;
	radius: number;
	color: string;
	moveX: number;
	moveY: number;
	speed: number;
	clientTime?: number;
}

interface Config {
	maxX: number;
	maxY: number;
}

const SECOND = 1000;
const FPS = 120;
const FPS_INTERVAL = Math.floor(SECOND / FPS);

const GameCanvas: React.FC = () => {
	const canvasRef = useRef<HTMLCanvasElement | null>(null);
	const [clientUid, setClientUid] = useState<string | null>(null);
	const [clientConfig, setClientConfig] = useState<Config | null>(null);
	const [clientPlayers, setClientPlayers] = useState<Record<string, Player>>({});
	const [moveX, setMoveX] = useState(0);
	const [moveY, setMoveY] = useState(0);

	const ws = useRef<WebSocket | null>(null);

	const setCanvasSize = useCallback(() => {
		if (canvasRef.current) {
			canvasRef.current.height = window.innerHeight;
			canvasRef.current.width = window.innerWidth;
		}
	}, []);

	useEffect(() => {
		setCanvasSize();
		window.addEventListener('resize', setCanvasSize);
		return () => window.removeEventListener('resize', setCanvasSize);
	}, [setCanvasSize]);

	const getRoundedInterval = () => {
		const now = Date.now();
		const value =  Math.floor(now / FPS_INTERVAL) * FPS_INTERVAL;
		return value;
	};

	const drawCircle = (context: CanvasRenderingContext2D, { x, y, radius, color }: Player) => {
		context.beginPath();
		context.arc(x, y, radius, 0, 2 * Math.PI, false);
		context.fillStyle = color;
		context.fill();

		context.lineWidth = Math.PI;

		context.strokeStyle = darkenColor(color, 30);
		context.stroke();
		context.lineWidth = 1;
	};

	function darkenColor(color: string, percent: number): string {
		// Validate the input color
		if (!/^#([0-9A-Fa-f]{6})$/.test(color)) {
			throw new Error("Invalid color format. Expected #RRGGBB.");
		}

		// Parse the red, green, and blue components from the color
		const r = parseInt(color.slice(1, 3), 16);
		const g = parseInt(color.slice(3, 5), 16);
		const b = parseInt(color.slice(5, 7), 16);

		// Function to darken a single channel
		const darkenChannel = (channel: number): number => 
			Math.max(0, Math.min(255, Math.floor(channel * (1 - percent / 100))));

		const newR = darkenChannel(r);
		const newG = darkenChannel(g);
		const newB = darkenChannel(b);

		// Convert the new RGB values back to a hex string
		return `#${newR.toString(16).padStart(2, "0")}${newG.toString(16).padStart(2, "0")}${newB.toString(16).padStart(2, "0")}`;
	}

	const drawGrid = (context: CanvasRenderingContext2D, playerX: number, playerY: number) => {
		const dx = 70;
		const dy = 70;

		const gridOffsetX = playerX % dx;
		const gridOffsetY = playerY % dy;

		context.beginPath();

		for (let x = -gridOffsetX; x <= context.canvas.width; x += dx) {
			context.moveTo(x, 0);
			context.lineTo(x, context.canvas.height);
		}

		for (let y = -gridOffsetY; y <= context.canvas.height; y += dy) {
			context.moveTo(0, y);
			context.lineTo(context.canvas.width, y);
		}

		context.strokeStyle = 'lightgray';
		context.stroke();
	};

	const drawDebugger = (context: CanvasRenderingContext2D, player: Player) => {
		const debuggableValues = [
			['x', Math.floor(player.x)],
			['y', Math.floor(player.y)],
			['players', Object.keys(clientPlayers).length],
			['latency', Date.now() - player.clientTime],
		];

		const fontSize = 24;
		const padding = 4;
		context.font = `normal ${fontSize}px Arial, sans-serif`;
		context.fillStyle = 'black';

		debuggableValues.forEach(([key, val], index) => {
			context.fillText(`${key}: ${val}`, fontSize, (index + 1) * (fontSize + padding));
		});
	};

	const draw = useCallback(() => {
		const canvas = canvasRef.current;
		if (!canvas) return;
		const context = canvas.getContext('2d');
		if (!context || !clientUid || !clientConfig) return;

		context.clearRect(0, 0, canvas.width, canvas.height);

		const clientPlayer = clientPlayers[clientUid];
		if (!clientPlayer) return;

		const xOffset = canvas.width / 2 - clientPlayer.x;
		const yOffset = canvas.height / 2 - clientPlayer.y;

		drawGrid(context, clientPlayer.x, clientPlayer.y);

		Object.keys(clientPlayers).forEach((uid) => {
			const player = clientPlayers[uid];
			drawCircle(context, {
				...player,
				x: player.x + xOffset,
				y: player.y + yOffset,
			});
		});

		drawDebugger(context, clientPlayer);
	}, [clientPlayers, clientUid, clientConfig]);

	const updatePlayer = useCallback(() => {
		if (!clientUid || !clientConfig) return;

		setClientPlayers((prevPlayers) => {
			const player = prevPlayers[clientUid];
			if (!player) return prevPlayers;

			const verifiedMoveX = Math.max(-1, Math.min(1, moveX));
			const verifiedMoveY = Math.max(-1, Math.min(1, moveY));

			const updatedPlayer = {
				...player,
				x: Math.max(0, Math.min(clientConfig.maxX, player.x + verifiedMoveX * player.speed)),
				y: Math.max(0, Math.min(clientConfig.maxY, player.y + verifiedMoveY * player.speed)),
				moveX: verifiedMoveX,
				moveY: verifiedMoveY,
			};

			ws.current?.send(
				JSON.stringify({
					uid: clientUid,
					moveX,
					moveY,
					clientTime: getRoundedInterval(),
				})
			);

			return {
				...prevPlayers,
				[clientUid]: updatedPlayer,
			};
		});
	}, [clientUid, clientConfig, moveX, moveY]);

	useEffect(() => {
		const interval = setInterval(updatePlayer, FPS_INTERVAL);
		return () => clearInterval(interval);
	}, [updatePlayer]);

	useEffect(() => {
		ws.current = new WebSocket('ws://localhost:8080/ws');

		ws.current.onmessage = (event) => {
			const data = JSON.parse(event.data);

			if (data.action === 'join') {
				setClientUid(data.uid);
				setClientConfig(data.config);
				setClientPlayers(data.players);
			}

			if (data.action === 'sync') {
				setClientPlayers(data.players);
			}
		};

		return () => {
			ws.current?.close();
		};
	}, []);

	useEffect(() => {
		const canvas = canvasRef.current;
		if (!canvas) return;

		const handleMouseMove = (event: MouseEvent) => {
			const rect = canvas.getBoundingClientRect();
			const mouseX = event.clientX - rect.left - canvas.width / 2;
			const mouseY = event.clientY - rect.top - canvas.height / 2;
			setMoveX(Math.min(1, Math.max(-1, mouseX / canvas.width)));
			setMoveY(Math.min(1, Math.max(-1, mouseY / canvas.height)));
		};

		canvas.addEventListener('mousemove', handleMouseMove);
		return () => canvas.removeEventListener('mousemove', handleMouseMove);
	}, []);

	useEffect(() => {
		const animationFrame = requestAnimationFrame(draw);
		return () => cancelAnimationFrame(animationFrame);
	}, [draw]);

	return <canvas ref={canvasRef} />;
};

export default GameCanvas;
