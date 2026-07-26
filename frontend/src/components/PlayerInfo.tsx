import React from "react";
import Image from 'next/image';

interface PlayerInfoProps {
  username: string;
  elo: number;
  isConnected?: boolean;
  isCurrentPlayer?: boolean;
}

const PlayerInfo: React.FC<PlayerInfoProps> = ({
  username,
  elo,
  isConnected = true,
  isCurrentPlayer = false
}) => {


  return (
    <div className={`flex items-center space-x-2 bg-white p-2 rounded ${!isConnected ? 'opacity-50' : ''}`}>

      {/* Avatar */}
      < Image
        alt="Player avatar"
        className="w-10 h-10 rounded"
        height="40"
        width="40"
        src="/assets/avatar.png"
      />

      {/* Player Info */}
      <div className="flex-1">
        <div className="flex items-center space-x-1 text-sm font-sans text-gray-900 font-semibold">

          {/* Username */}
          <span className={isCurrentPlayer ? 'text-blue-600' : ''}>
            {username}
          </span>

          {/* ELO Rating */}
          <span className="font-normal text-gray-600">
            ({elo})
          </span>


          {/* Connection Status */}
          <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></span>

        </div>



      </div>

    </div>
  );
};

export default PlayerInfo;