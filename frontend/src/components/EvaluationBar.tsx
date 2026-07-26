'use client';

import React from 'react';

interface EvaluationBarProps {
  evaluation: number; // -10 to +10 (negative = black advantage, positive = white advantage)
  playerColor: 'white' | 'black';
  className?: string;
}

const EvaluationBar: React.FC<EvaluationBarProps> = ({
  evaluation,
  playerColor,
  className = ''
}) => {

  const clampedEval = Math.max(-10, Math.min(10, evaluation));

  const whiteAdvantagePercentage = ((clampedEval + 10) / 20) * 100;
  
  let whiteSectionPercentage;
  
  if (playerColor === 'white') {
    whiteSectionPercentage = whiteAdvantagePercentage;
  } else {
    whiteSectionPercentage = 100 - whiteAdvantagePercentage;
  }
  
  // Format evaluation display
  const getEvalDisplay = () => {
    if (Math.abs(clampedEval) < 0.1) return '0.0';
    if (clampedEval > 0) return `+${clampedEval.toFixed(1)}`;
    return clampedEval.toFixed(1);
  };



  return (
    <div className={`flex flex-col h-full ${className}`}>


      <div className="relative flex-1 bg-gray-900 overflow-hidden border border-gray-700 w-8 mx-auto">

        <div 
          className="absolute bottom-0 left-0 right-0 bg-white transition-all duration-500 ease-out"
          style={{ 
            height: `${whiteSectionPercentage}%`,
          }}
        ></div>

      </div>

      {/* Evaluation display */}
      <div className="mt-3 text-center">
        <div className="text-lg font-bold text-gray-900">
          {getEvalDisplay()}
        </div>

      </div>
    </div>
  );
};

export default EvaluationBar;