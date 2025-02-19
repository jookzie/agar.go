import React, { useEffect, useRef, useState, useCallback } from 'react';
import clone from 'just-clone';

interface Player {
	name: string;
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
const FPS = 500;
const FPS_INTERVAL = Math.floor(SECOND / FPS);

const GameCanvas: React.FC = () => {
	const canvasRef = useRef<HTMLCanvasElement | null>(null);
	const [clientUid, setClientUid] = useState<string | null>(null);
	const [clientConfig, setClientConfig] = useState<Config | null>(null);
	const [clientPlayers, setClientPlayers] = useState<Record<string, Player>>({});
	const [clientFeedmap, setClientFeedmap] = useState<Array<Array<number>>>([]);
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
		const value = Math.floor(now / FPS_INTERVAL) * FPS_INTERVAL;
		return value;
	};

	const drawCircle = (context: CanvasRenderingContext2D, { x, y, radius, name, color }: Player) => {
		context.beginPath();
		context.arc(x, y, radius, 0, 2 * Math.PI, false);
		context.fillStyle = color || "black";
		context.fill();

		context.lineWidth = 2 * Math.PI;
		context.strokeStyle = darkenColor(color || "black", 30);
		context.stroke();

		const text = name;
		const coord = [
			x-13*text.length/2, 
			y-50+Math.min(50, radius)
		];
		context.font = '25px Sans-serif';
		context.strokeStyle = 'black';
		context.lineWidth = 5;
		context.strokeText(text, coord[0], coord[1]);
		context.fillStyle = 'white';
		context.fillText(text, coord[0], coord[1]);

		context.lineWidth = 1;
	};

	const drawFeedCircle = (context: CanvasRenderingContext2D, x: number, y: number, xOffset: number, yOffset: number) => {
		context.beginPath();
		context.arc(x + xOffset, y + yOffset, (x + y) % 3 + 10, 0, 2 * Math.PI, false);
		context.fillStyle = '#' + (((x + y) * 1234567) & 0xFFFFFF).toString(16).padStart(6, '0')
		context.fill();
	};

	const drawLeaderboard = (ctx: CanvasRenderingContext2D) => {
		const canvas = canvasRef.current;
		if (!canvas) return;
		 // Sort players by radius in descending order
		const clonedPlayers = clone(Object.values(clientPlayers))
		const sortedPlayers = [...clonedPlayers].sort((a, b) => b.radius - a.radius);

		// Draw the leaderboard box
		const boxWidth = 200;
		const boxHeight = sortedPlayers.length * 30 + 40; // Dynamic height based on players
		const boxX = canvas.width - boxWidth - 10; // Top-right corner
		const boxY = 10;

		ctx.globalAlpha = 0.8; // Set opacity
		ctx.fillStyle = "black";
		ctx.fillRect(boxX, boxY, boxWidth, boxHeight);

		ctx.globalAlpha = 1.0; // Reset opacity
		ctx.fillStyle = "white";
		ctx.font = "16px Arial";

		// Draw the leaderboard title
		const title = "Leaderboard";
		const titleWidth = ctx.measureText(title).width;
		const titleX = boxX + (boxWidth - titleWidth) / 2; // Center the title
		ctx.fillText(title, titleX, boxY + 25);

		// Draw player rankings
		sortedPlayers.forEach((player, index) => {
		  const text = `${index + 1}. ${player.name}`;
		  ctx.fillText(text, boxX + 10, boxY + 50 + index * 30); // Adjust Y position for the title
		});
	};

	function darkenColor(color: string, percent: number): string {
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
		
		const latencies = [];

		const latency = Math.round((Date.now() - (player.clientTime || 0)) / 5) * 5

		const debuggableValues = [
			['x', Math.floor(player.x)],
			['y', Math.floor(player.y)],
			['players', Object.keys(clientPlayers).length],
			['latency', `${String(latency).padStart(2, '0')} ms`],
			['feed points', clientFeedmap.length],
			['score', Math.floor(player.radius * 10 - 200)],
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
		if (!context || !clientUid || !clientConfig || !clientFeedmap) return;

		context.clearRect(0, 0, canvas.width, canvas.height);

		const player = clientPlayers[clientUid];
		if (!player) return;

		const xOffset = canvas.width / 2 - player.x;
		const yOffset = canvas.height / 2 - player.y;

		drawGrid(context, player.x, player.y);

		clientFeedmap.forEach((point: Array<number>) => {
			drawFeedCircle(
				context,
				point[0], 
				point[1],
				xOffset,
				yOffset,
			)
		});

		const clonedPlayers = clone(Object.values(clientPlayers))
		clonedPlayers
			.sort((a, b) => a.radius - b.radius)
			.forEach((p) => {
				drawCircle(context, {
					...p,
					x: p.x + xOffset,
					y: p.y + yOffset,
				});
			});

		if (import.meta.env.DEV) {
			drawDebugger(context, player)
		}
		drawLeaderboard(context);
	}, [clientPlayers, clientUid, clientConfig, clientFeedmap]);

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
		ws.current = new WebSocket(import.meta.env.VITE_BACKEND_URL);

		ws.current.onmessage = (event) => {
			const data = JSON.parse(event.data);

			if (data.action === 'join') {
				setClientUid(data.uid);
				setClientConfig(data.config);
				setClientPlayers(data.players);
				setClientFeedmap(data.feedmap);
			}

			if (data.action === 'sync') {
				setClientPlayers(prev => {
					prev[data.uid] = data.player;
					if (data.eatenPlayer) {
						delete prev[data.eatenPlayer];
					}
					return prev;
				});
				data.eatenPoint && setClientFeedmap((prev: Array<Array<number>>) => {
					const idx = prev.findIndex(item => item[0] == data.eatenPoint[0] && item[1] == data.eatenPoint[1]);
					if (idx == -1) throw Error("Eaten point was not found in feed map.")

					let feedmap = structuredClone(prev);
					feedmap.splice(idx, 1);
					feedmap = [...feedmap, ...data.addedPoints || []];
					return feedmap;
				});
			}

			if(data.action == "delete") {
				setClientPlayers(prev => {
					delete prev[data.uid];
					return prev;
				});
			}
		};

		ws.current.onclose = (event) => {
			window.location.reload()
		};


		return () => {
			ws.current?.close();
		};
	}, []);

	const updateMovement = useCallback(() => {
		const canvas = canvasRef.current;
		if (!canvas) return;

		// Get the canvas dimensions and calculate movement offsets
		const rect = canvas.getBoundingClientRect();

		const handleMouseMove = (event: MouseEvent) => {
			const rect = canvas.getBoundingClientRect();

			// Calculate mouse position relative to the center of the canvas
			const mouseX = event.clientX - rect.left - canvas.width / 2;
			const mouseY = event.clientY - rect.top - canvas.height / 2;

			// Calculate the magnitude (distance from the center)
			const magnitude = Math.sqrt(mouseX ** 2 + mouseY ** 2);

			// Normalize the direction vector
			const directionX = magnitude === 0 ? 0 : mouseX / magnitude;
			const directionY = magnitude === 0 ? 0 : mouseY / magnitude;

			const verifiedMoveX = Math.max(-1, Math.min(1, directionX));
			const verifiedMoveY = Math.max(-1, Math.min(1, directionY));

			// Set the movement values
			setMoveX(verifiedMoveX);
			setMoveY(verifiedMoveY);
		};

		// Add the mousemove event listener
		canvas.addEventListener("mousemove", handleMouseMove);

		// Cleanup function to remove the event listener
		return () => {
			canvas.removeEventListener("mousemove", handleMouseMove);
		};
	}, []);


	useEffect(() => {
		const animationFrame = requestAnimationFrame(draw);
		return () => cancelAnimationFrame(animationFrame);
	}, [draw]);

	useEffect(() => {
		const animationFrame = requestAnimationFrame(updateMovement);
		return () => cancelAnimationFrame(animationFrame);
	}, [updateMovement]);

	return <canvas ref={canvasRef} />;
};

export default GameCanvas;
