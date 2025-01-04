import './styles.scss';
import { useState, useEffect, useRef } from "react";

export default function MainPage() {
	const [messages, setMessages] = useState<string[]>([]);
	const [input, setInput] = useState<string>("");
	const ws = useRef<WebSocket | null>(null);

	useEffect(() => {
		ws.current = new WebSocket("ws://localhost:8080/ws");

		ws.current.onmessage = (event: MessageEvent) => {
			setMessages((prev) => [...prev, event.data]);
		};

		return () => {
			ws.current?.close();
		};
	}, []);

	const sendMessage = () => {
		if (ws.current && input) {
			ws.current.send(JSON.stringify(input));
			setInput("");
		}
	};

	return (
		<div className="chatroom">
			<div className="messages">
				{messages.map((msg: string, idx: number) => (
					<div key={idx} className="message">
						{msg}
					</div>
				))}
			</div>
			<div className="input">
				<input
					type="text"
					value={input}
					onChange={(e) => setInput(e.target.value)}
					placeholder="Type a message..."
				/>
				<button onClick={sendMessage}>Send</button>
			</div>
		</div>
	);
}
