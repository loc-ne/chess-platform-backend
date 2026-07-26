import React from 'react';
import Image from 'next/image';

interface ReplayControlsProps {
  onStart?: () => void;
  onPrevious?: () => void;
  onNext?: () => void;
  onEnd?: () => void;
  onAutoReplay?: () => void;
  isReplaying: boolean;
  disabledStart?: boolean;
  disabledPrevious?: boolean;
  disabledNext?: boolean;
  disabledEnd?: boolean;
  currentMove?: number;
  totalMoves?: number;
}

const iconClass = "w-4 h-4";
const iconSize = 16;

const ReplayControls: React.FC<ReplayControlsProps> = ({
  onStart,
  onPrevious,
  onNext,
  onEnd,
  onAutoReplay,
  isReplaying,
  disabledStart,
  disabledPrevious,
  disabledNext,
  disabledEnd,
}) => (
  <div className="flex gap-3 justify-center items-center mt-2">
    <button
      onClick={onStart}
      disabled={disabledStart}
      className="p-2 rounded hover:bg-gray-100 disabled:opacity-50"
      aria-label="Về đầu"
    >
      < Image src="/assets/start-icon.png" alt="Start" className={iconClass} width={iconSize}
        height={iconSize} />
    </button>
    <button
      onClick={onPrevious}
      disabled={disabledPrevious}
      className="p-2 rounded hover:bg-gray-100 disabled:opacity-50"
      aria-label="Lùi lại"
    >
      < Image src="/assets/previous-icon.png" alt="Previous" className={iconClass} width={iconSize}
        height={iconSize} />
    </button>
    <button
      onClick={onAutoReplay}
      className="p-2 rounded hover:bg-gray-100 disabled:opacity-50"
      aria-label={isReplaying ? "Tạm dừng" : "Tự động chạy"}
    >
      <Image
        src={isReplaying ? "/assets/pause.png" : "/assets/play-icon.png"}
        alt={isReplaying ? "Pause" : "Play"}
        className={iconClass}
        width={iconSize}
        height={iconSize}
      />
    </button>


    <button
      onClick={onNext}
      disabled={disabledNext}
      className="p-2 rounded hover:bg-gray-100 disabled:opacity-50"
      aria-label="Tiến tới"
    >
      < Image src="/assets/next-icon.png" alt="Next" className={iconClass} width={iconSize}
        height={iconSize} />
    </button>
    <button
      onClick={onEnd}
      disabled={disabledEnd}
      className="p-2 rounded hover:bg-gray-100 disabled:opacity-50"
      aria-label="Về cuối"
    >
      < Image src="/assets/end-icon.png" alt="End" className={iconClass} width={iconSize}
        height={iconSize} />
    </button>
  </div>
);

export default ReplayControls;