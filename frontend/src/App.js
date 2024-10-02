import React, { useState, useEffect } from 'react';
import Board from './components/Board';
import { GoogleLogin } from '@react-oauth/google';
import { startGame, makeMove } from './api';
import axios from 'axios';
function App() {

  const [game, setGame] = useState(null);
  const [status, setStatus] = useState('');
  const [name, setName] = useState("");
  const [userID, setuserID] = useState("");

  // Start a new game when the component mounts
  useEffect(() => {
    getQueryParams();
  }, []);

  // Handle player moves
  const handleMove = async (position) => {
    if (game && game.board[position] === '' && game.status === 'ongoing') {
      const updatedGame = await makeMove(game.id, position);
      setGame(updatedGame);
      setStatus(updatedGame.status);
      console.log(updatedGame.status);
      if (updatedGame.status != "ongoing") {
        const googleLogoutUrl = 'https://accounts.google.com/Logout';
        window.location.href = googleLogoutUrl;
        window.location.href = 'http://192.168.1.151:1975';
      }
    }
  };
  async function initGame(userID) {
    const newGame = await startGame(userID);
    setGame(newGame);
    setStatus('ongoing');
  }
  function getQueryParams() {
    var url_string = window.location;
    if (url_string.search) {
      var url = decodeURIComponent(url_string.search.slice(1, 22))
      var name1 = decodeURIComponent(url_string.search.slice(22))
      console.log(url, name1)
      setuserID(url)
      setName(name1)
      initGame(url);
    }
    return 1;
  }

  const handleLogin = () => {
    window.location.href = 'http://192.168.1.151:1978/auth/google/login';
    getQueryParams();
  };
  const logOut = () => {
    const googleLogoutUrl = 'https://accounts.google.com/Logout';
    window.location.href = googleLogoutUrl;
    window.location.href = 'http://192.168.1.151:1975';
  }




  // Render the board and game status
  return (
    <div className="App">
      <h1>Tic-Tac-Toe</h1>
      <button onClick={handleLogin} style={{ display: (game) ? "none" : "block" }}>
        Login with Google
      </button>
      <button onClick={logOut} style={{ display: (game && game.Flag) ? "block" : "none" }}>
        Logout
      </button>
      <p>{name ? "Welcome to " + name + "!!!!!!!!" : ""}</p>
      {(game && game.Flag) ? (
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
