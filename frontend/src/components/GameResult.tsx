'use client';

import React, { useState } from 'react';

// Types
interface GameResultData {
    result: string;
    reason: string;
}

interface ChessResultDialogProps {
    isOpen?: boolean;
    gameResult?: GameResultData;
    onRematch?: () => void;
    onNewGame?: () => void;
    onClose?: () => void;
}

const GameResultDialog: React.FC<ChessResultDialogProps> = ({
    isOpen: propIsOpen,
    gameResult: propGameResult,
    onRematch,
    onNewGame,
    onClose
}) => {
    const [internalIsOpen, setInternalIsOpen] = useState<boolean>(true);
    const [internalGameResult] = useState<GameResultData>({
        result: 'white-wins',
        reason: 'Chiếu hết'
    });

    const isOpen = propIsOpen !== undefined ? propIsOpen : internalIsOpen;
    const gameResult = propGameResult || internalGameResult;

    const handleRematch = (): void => {
        if (onRematch) {
            onRematch();
        } else {
            setInternalIsOpen(false);
        }
    };

    const handleNewGame = (): void => {
        if (onNewGame) {
            onNewGame();
        } else {
            setInternalIsOpen(false);
        }
    };

    const handleClose = (): void => {
        if (onClose) {
            onClose();
        } else {
            setInternalIsOpen(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="relative bg-white rounded-xl shadow-lg max-w-sm w-full mx-4 p-6 z-10">
                <button
                    onClick={handleClose}
                    className="absolute top-3 right-4 text-gray-400 hover:text-gray-700 text-2xl font-bold"
                    aria-label="Đóng"
                >
                    ×
                </button>
                <div className="text-center mb-4">
                    <div className="text-lg font-semibold text-blue-700 mb-1">{gameResult.result}</div>
                    <div className="text-gray-600 text-base mb-2">Lý do: {gameResult.reason}</div>
                </div>
                <div className="flex flex-col items-center gap-3 mt-6">
                    <div className="flex gap-3 w-full max-w-2xl">
                        <button
                            onClick={handleRematch}
                            className="flex-1 py-2.5 px-6 text-sm font-semibold text-gray-800 bg-white border-2 border-gray-400 rounded-lg shadow-sm hover:bg-gray-100 hover:text-blue-700 hover:border-blue-500 focus:outline-none focus:ring-4 focus:ring-blue-200 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700 transition-all duration-200 ease-in-out"
                        >
                            Tái đấu
                        </button>
                        <button
                            onClick={handleNewGame}
                            className="flex-1 py-2.5 px-6 text-sm font-semibold text-gray-800 bg-white border-2 border-gray-400 rounded-lg shadow-sm hover:bg-gray-100 hover:text-blue-700 hover:border-blue-500 focus:outline-none focus:ring-4 focus:ring-blue-200 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700 transition-all duration-200 ease-in-out"
                        >
                            Trận mới
                        </button>
                    </div>

                    <div className="w-full max-w-2xl">
                        <button
                            onClick={() => alert('Phân tích ván đấu!')}
                            className="w-full py-2.5 px-6 text-sm font-semibold text-gray-800 bg-white border-2 border-gray-400 rounded-lg shadow-sm hover:bg-gray-100 hover:text-blue-700 hover:border-blue-500 focus:outline-none focus:ring-4 focus:ring-blue-200 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700 transition-all duration-200 ease-in-out"
                        >
                            Phân tích ván đấu
                        </button>
                    </div>
                </div>

            </div>
        </div>
    );
};

export default GameResultDialog;