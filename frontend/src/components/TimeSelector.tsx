'use client';

import React, { useState } from 'react';

// Interfaces
interface TimeCategory {
  name: string;
  times: string[];
  icon: string;
  color: string;
  description: string;
}

interface TimeSelectorProps {
  initialTime?: string;
  initialCategory?: string;
  onTimeSelect?: (time: string, category: string) => void;
  disabled?: boolean;
  size?: 'small' | 'medium' | 'large';
  className?: string;
}

const TimeSelector: React.FC<TimeSelectorProps> = ({
  initialTime = '5 + 0',
  initialCategory = 'Blitz',
  onTimeSelect,
  disabled = false,
  size = 'medium',
  className = '',
}) => {
  const [selectedTime, setSelectedTime] = useState(initialTime);
  const [selectedCategory, setSelectedCategory] = useState(initialCategory);
  const [activeCategory, setActiveCategory] = useState(initialCategory);

  const timeCategories: TimeCategory[] = [
    { 
      name: 'Bullet', 
      times: ['1 + 0', '1 + 1', '2 + 1'], 
      icon: '⚡',
      color: 'from-red-600 to-orange-600',
      description: 'Lightning fast games'
    },
    { 
      name: 'Blitz', 
      times: ['3 + 0', '3 + 2', '5 + 0'], 
      icon: '🔥',
      color: 'from-orange-600 to-amber-600',
      description: 'Quick tactical battles'
    },
    { 
      name: 'Rapid', 
      times: ['10 + 0', '10 + 5', '15 + 10'], 
      icon: '⏰',
      color: 'from-blue-600 to-cyan-600',
      description: 'Balanced thinking time'
    },
    { 
      name: 'Classical', 
      times: ['30 + 0', '30 + 20'], 
      icon: '🎯',
      color: 'from-purple-600 to-pink-600',
      description: 'Deep strategic play'
    },
  ];

  const getCurrentCategory = () => {
    return timeCategories.find(cat => cat.name === selectedCategory) || timeCategories[1];
  };

  const handleTimeSelect = (time: string, category: string) => {
    setSelectedTime(time);
    setSelectedCategory(category);
    onTimeSelect?.(time, category);
  };

  const currentCat = getCurrentCategory();

  return (
    <div className={`w-full ${className}`}>
      {/* Selected Time Display */}
      <div className="mb-6">
        <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-lg">
          <div className="text-center">
            <div className="flex items-center justify-center gap-3 mb-3">
              <span className="text-3xl">{currentCat.icon}</span>
              <div>
                <div className="text-2xl font-bold text-gray-900">{selectedTime}</div>
                <div className="text-sm text-gray-600">{selectedCategory}</div>
              </div>
            </div>
            <div className={`text-xs text-gray-500 bg-gradient-to-r ${currentCat.color} bg-clip-text text-transparent font-medium`}>
              {currentCat.description}
            </div>
          </div>
        </div>
      </div>

      {/* Category Tabs */}
      <div className="grid grid-cols-2 gap-2 mb-6">
        {timeCategories.map((category) => {
          const isActive = activeCategory === category.name;
          return (
            <button
              key={category.name}
              className={`
                relative p-4 rounded-xl font-medium transition-all duration-300 border
                ${isActive 
                  ? 'bg-gray-100 border-gray-300 text-gray-900 shadow-md scale-105' 
                  : 'bg-white border-gray-200 text-gray-600 hover:bg-gray-50 hover:border-gray-300'
                }
                ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
              `}
              onClick={() => !disabled && setActiveCategory(category.name)}
              disabled={disabled}
            >
              <div className="flex items-center gap-2 justify-center">
                <span className="text-lg">{category.icon}</span>
                <span className="text-sm font-semibold">{category.name}</span>
              </div>
              
              {/* Active indicator */}
              {isActive && (
                <div className={`absolute inset-0 rounded-xl bg-gradient-to-r ${category.color} opacity-5`}></div>
              )}
            </button>
          );
        })}
      </div>

      {/* Time Options for Active Category */}
      <div className="space-y-4">
        {timeCategories
          .filter(cat => cat.name === activeCategory)
          .map((category) => (
            <div key={category.name} className="space-y-3">
              {/* Category Description */}
              <div className="flex items-center gap-3 mb-4">
                <div className={`w-1 h-6 bg-gradient-to-b ${category.color} rounded-full`}></div>
                <div>
                  <h3 className="text-gray-900 font-semibold flex items-center gap-2">
                    <span className="text-xl">{category.icon}</span>
                    {category.name}
                  </h3>
                  <p className="text-xs text-gray-500">{category.description}</p>
                </div>
              </div>

              {/* Time Buttons Grid */}
              <div className="grid grid-cols-1 gap-3">
                {category.times.map((time) => {
                  const isSelected = selectedTime === time && selectedCategory === category.name;
                  return (
                    <button
                      key={time}
                      className={`
                        relative p-4 rounded-xl font-semibold transition-all duration-300 border group overflow-hidden
                        ${isSelected
                          ? 'bg-gradient-to-r from-green-50 to-emerald-50 border-green-300 text-gray-900 shadow-lg scale-105'
                          : 'bg-white border-gray-200 text-gray-700 hover:bg-gray-50 hover:border-gray-300 hover:scale-102'
                        }
                        ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
                      `}
                      onClick={() => !disabled && handleTimeSelect(time, category.name)}
                      disabled={disabled}
                    >
                      <div className="relative z-10 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className={`w-3 h-3 rounded-full transition-all duration-300 ${
                            isSelected ? 'bg-green-500 shadow-lg shadow-green-500/30' : 'bg-gray-400'
                          }`}></div>
                          <span className="text-lg">{time}</span>
                        </div>
                        
                        {/* Time format explanation */}
                        <div className="text-xs text-gray-500">
                          {time.includes('+') ? 'minutes + increment' : 'minutes'}
                        </div>
                      </div>

                      {/* Selection indicator */}


                      {/* Hover effect */}
                      <div className="absolute inset-0 bg-gradient-to-r from-gray-50/0 to-gray-50/80 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
      </div>

  
    </div>
  );
};

export default TimeSelector;