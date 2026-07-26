import React, { useEffect } from 'react';

interface ClockProps {
  initialMinutes: number;
  initialSeconds?: number;
  increment?: number;
  isActive?: boolean;
  onTimeUp?: () => void;
}

const Clock = ({ 
  initialMinutes, 
  initialSeconds = 0, 
  isActive = false,
  onTimeUp 
}: ClockProps) => {
  const [minutes, setMinutes] = React.useState(initialMinutes);
  const [seconds, setSeconds] = React.useState(initialSeconds);
  const timerRef = React.useRef<number | null>(null);

useEffect(() => {
    if (!isActive) {
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
      return;
    }

    timerRef.current = window.setInterval(() => {
      const totalSeconds = minutes * 60 + seconds;
      
      if (totalSeconds <= 0) {
        onTimeUp?.();
        return;
      }

      const newTotalSeconds = totalSeconds - 1;
      const newMinutes = Math.floor(newTotalSeconds / 60);
      const newSeconds = newTotalSeconds % 60;

      setMinutes(newMinutes);
      setSeconds(newSeconds);
    }, 1000);

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    };
  }, [isActive, onTimeUp, minutes, seconds]);

  const formatTime = (min: number, sec: number) => {
    return `${min.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}`;
  };

  const isLowTime = minutes < 1 && seconds <= 30;
  const isVeryLowTime = minutes === 0 && seconds <= 10;

  return (
    <div 
      className={`
        text-2xl font-mono font-bold p-3 
        ${isVeryLowTime ? 'bg-red-600 text-black animate-pulse' : 
          isLowTime ? 'bg-yellow-500 text-black' : 
          'bg-white text-black'}
      `}
    >
       {formatTime(minutes, seconds)}

    </div>
  );
};

export default Clock;