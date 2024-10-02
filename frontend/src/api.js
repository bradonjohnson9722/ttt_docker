import axios from 'axios';

const api = "http://192.168.1.151:1972";
// Start a new game by making a POST request to the backend
export const startGame = async (userID) => {
  const response = await axios.post(`${api}/start-game?userid=${userID}`);
  return response.data;
};

// Make a move by sending the player's move to the backend
export const makeMove = async (gameId, position) => {
  const response = await axios.post(`${api}/make-move?id=${gameId}`, {
    player: "user",
    position,
  });
  return response.data;
};
