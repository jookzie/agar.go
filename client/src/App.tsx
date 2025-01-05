import GameCanvas from './GameCanvas.tsx'

export default function App() {
	console.log(import.meta.env.VITE_BACKEND_URL);
	return (
		<GameCanvas />
	);
}
