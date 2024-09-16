import React, { useState, useEffect } from 'react';
import Board from './components/Board';
import { startGame, makeMove } from './api';

function App() {
  const [game, setGame] = useState(null);
  const [status, setStatus] = useState('');

  // Start a new game when the component mounts
  useEffect(() => {
    async function initGame() {
      const newGame = await startGame();
      setGame(newGame);
      setStatus('Game ongoing');
    }
    initGame();
  }, []);

  // Handle player moves
  const handleMove = async (position) => {
    if (game && game.board[position] === '' && game.status === 'ongoing') {
      const updatedGame = await makeMove(game.id, position);
      setGame(updatedGame);
      setStatus(updatedGame.status);
    }
  };

  // Render the board and game status
  return (
    <div className="App">
      <h1>Tic-Tac-Toe</h1>
      {game ? (
        <>
          <Board board={game.board} onClick={handleMove} />
          <p>Status: {status}</p>
        </>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
}

export default App;
