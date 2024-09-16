import React from 'react';
import Square from './Square';

function Board({ board, onClick }) {
  return (
    <div className="board">
      {board.map((value, index) => (
        <Square key={index} value={value} onClick={() => onClick(index)} />
      ))}
    </div>
  );
}

export default Board;
