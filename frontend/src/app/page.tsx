'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import Image from 'next/image';

export default function Home() {
  const router = useRouter();
  const { user } = useAuth();

  const handlePlayOnline = () => {
    if (user) {
      router.push('/play/online');
    } else {
      router.push('/login');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Header */}

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h2 className="text-4xl font-bold text-gray-900 mb-4">
            Chơi Cờ Vua Online
          </h2>
        </div>

        {/* Play Options */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          {/* Quick Play */}
          <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200 hover:shadow-xl transition-shadow">
            <div className="text-center">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" id="play-button">
                  <path fill="#24c4d8" d="M6.418 27.157V4.843a1 1 0 0 1 1.544-.839l17.166 11.158a1 1 0 0 1 0 1.676L7.962 27.996a1 1 0 0 1-1.544-.839Z"></path>
                  <path fill="#02abbf" d="M25.582 16a.992.992 0 0 1-.455.838L7.963 27.996a1 1 0 0 1-1.545-.839V16Z"></path>
                  <path d="M7.418 29.157a2 2 0 0 1-2-2V4.843a2 2 0 0 1 3.09-1.677l17.164 11.157a2 2 0 0 1 0 3.354L8.507 28.834a1.998 1.998 0 0 1-1.09.323Zm0-24.315v22.32L24.582 16Z"></path>
                </svg>


              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">Chơi Nhanh</h3>
              <p className="text-gray-600 mb-4">Tìm đối thủ ngay lập tức</p>
              <button
                onClick={handlePlayOnline}
                className="w-full px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 font-medium transition-colors"
              >
                Chơi ngay
              </button>
            </div>
          </div>

          {/* Time Controls */}
          <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200 hover:shadow-xl transition-shadow">
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" id="thinking">
                  <path fill="#fdc475" d="M16.51,31.5H27.45c-.69-2.09-1.52-4.59-2.18-6.54a5.66549,5.66549,0,0,1,.77-5.2A11.85615,11.85615,0,1,0,4.74,12.6c0,.2.01.4.02.61a8.33315,8.33315,0,0,1-.44,3.09L3.25,19.58a1.13011,1.13011,0,0,0,1.07,1.48h2.4v3.11a4.511,4.511,0,0,0,4.51,4.51h2.91a2.26618,2.26618,0,0,1,2.23,1.92Z"></path>
                  <path fill="#faa827" d="M7.72 24.17V21.06h-1v3.11a4.51106 4.51106 0 0 0 4.51 4.51h1A4.51106 4.51106 0 0 1 7.72 24.17zM4.25 19.58L5.32 16.3a8.33285 8.33285 0 0 0 .44-3.09c-.01-.21-.02-.41-.02-.61A11.84923 11.84923 0 0 1 16.94843.76477 11.8406 11.8406 0 0 0 4.74 12.6c0 .2.01.4.02.61a8.33285 8.33285 0 0 1-.44 3.09L3.25 19.58a1.13013 1.13013 0 0 0 1.07 1.48h1A1.13013 1.13013 0 0 1 4.25 19.58zM15.14 28.68h-1a2.26616 2.26616 0 0 1 2.23 1.92l.14.9h1l-.14-.9A2.26616 2.26616 0 0 0 15.14 28.68z"></path>
                  <path fill="#e5303e" d="M16.5,19a8,8,0,1,1,8-8A8.00917,8.00917,0,0,1,16.5,19Zm0-15a7,7,0,1,0,7,7A7.00849,7.00849,0,0,0,16.5,4Z"></path>
                  <path fill="#fff" d="M23.5,11a4.067,4.067,0,0,1-.03.5,6.83164,6.83164,0,0,1-1.68,4.08,5.94555,5.94555,0,0,1-.71.71A6.832,6.832,0,0,1,17,17.97a4.17762,4.17762,0,0,1-1,0,6.832,6.832,0,0,1-4.08-1.68,5.94555,5.94555,0,0,1-.71-.71A6.83164,6.83164,0,0,1,9.53,11.5a4.17762,4.17762,0,0,1,0-1,6.83181,6.83181,0,0,1,1.68-4.08l.71-.71A6.83272,6.83272,0,0,1,16,4.03a4.18606,4.18606,0,0,1,1,0,6.83272,6.83272,0,0,1,4.08,1.68,5.95007,5.95007,0,0,1,.71.71,6.83181,6.83181,0,0,1,1.68,4.08A4.06556,4.06556,0,0,1,23.5,11Z"></path>
                  <path fill="#e0e0e2" d="M12.21,15.58a6.83164,6.83164,0,0,1-1.68-4.08,4.17762,4.17762,0,0,1,0-1,6.83181,6.83181,0,0,1,1.68-4.08l.71-.71A6.83272,6.83272,0,0,1,17,4.03a4.18606,4.18606,0,0,0-1,0,6.83272,6.83272,0,0,0-4.08,1.68l-.71.71A6.83181,6.83181,0,0,0,9.53,10.5a4.17762,4.17762,0,0,0,0,1,6.83164,6.83164,0,0,0,1.68,4.08,5.94555,5.94555,0,0,0,.71.71A6.832,6.832,0,0,0,16,17.97a4.17762,4.17762,0,0,0,1,0,6.832,6.832,0,0,1-4.08-1.68A5.94555,5.94555,0,0,1,12.21,15.58Z"></path>
                  <path fill="#3c3b41" d="M16.5 10.5A.49971.49971 0 0 1 16 10V7.5a.5.5 0 0 1 1 0V10A.49971.49971 0 0 1 16.5 10.5zM19.43749 14.4375A.49838.49838 0 0 1 19.084 14.291l-2.22754-2.22754a.49995.49995 0 0 1 .707-.707L19.791 13.584a.5.5 0 0 1-.35352.85352zM17 3.83V4.5a.5.5 0 0 1-1 0V3.83a4.18606 4.18606 0 0 1 1 0zM17 17.5v.88664a4.17762 4.17762 0 0 1-1 0V17.5a.5.5 0 0 1 1 0zM11.707 5.6082l.47378.47378a.5.5 0 0 1-.70711.70711l-.47378-.47378A4.18606 4.18606 0 0 1 11.707 5.6082zM21.37314 15.27437l.62695.62695a4.17762 4.17762 0 0 1-.70711.70711l-.62695-.62695a.5.5 0 1 1 .70711-.70711zM9.22163 10.60833h.67a.5.5 0 1 1 0 1h-.67a4.18606 4.18606 0 0 1 0-1zM22.89167 10.60833h.88664a4.17762 4.17762 0 0 1 0 1h-.88664a.5.5 0 0 1 0-1zM10.99987 15.90136l.47378-.47378a.5.5 0 0 1 .70711.70711l-.47378.47378a4.18606 4.18606 0 0 1-.70711-.70711zM20.666 6.23519l.62695-.62695a4.17762 4.17762 0 0 1 .70711.70711l-.62695.62695a.5.5 0 1 1-.70711-.70711z"></path>
                  <circle cx="16.5" cy="11" r="1" fill="#e5303e"></circle>
                  <path fill="#e5303e" d="M16.5,3a8,8,0,1,0,8,8A8.01062,8.01062,0,0,0,16.5,3Zm6.92,9.01c-.01.04-.01.08-.02.12a2.88033,2.88033,0,0,1-.1.51,5.43057,5.43057,0,0,1-.15.54,6.666,6.666,0,0,1-.91,1.83c-.09.12-.18.25-.28.36-.05.07-.11.14-.17.21a5.94555,5.94555,0,0,1-.71.71c-.07.06-.14.12-.21.17-.11.1-.24.19-.36.28a6.66239,6.66239,0,0,1-1.83.91,5.431,5.431,0,0,1-.54.15,2.87477,2.87477,0,0,1-.51.1c-.04.01-.08.01-.12.02a2.93455,2.93455,0,0,1-.51.05,4.17762,4.17762,0,0,1-1,0,2.93455,2.93455,0,0,1-.51-.05c-.04-.01-.08-.01-.12-.02a2.87477,2.87477,0,0,1-.51-.1,5.431,5.431,0,0,1-.54-.15,6.66239,6.66239,0,0,1-1.83-.91c-.12-.09-.25-.18-.36-.28-.07-.05-.14-.11-.21-.17a5.94555,5.94555,0,0,1-.71-.71c-.06-.06995-.12-.14-.17-.21-.1-.11-.19-.24-.28-.36a6.666,6.666,0,0,1-.91-1.83,5.43057,5.43057,0,0,1-.15-.54,2.88033,2.88033,0,0,1-.1-.51c-.00995-.04-.00995-.08-.02-.12a2.93513,2.93513,0,0,1-.05-.51,4.17762,4.17762,0,0,1,0-1,2.934,2.934,0,0,1,.05-.51c.01-.04.01-.08.02-.12a2.87925,2.87925,0,0,1,.1-.51,5.43057,5.43057,0,0,1,.15-.54,6.66426,6.66426,0,0,1,.91-1.83c.09-.12.18-.25.28-.36.05-.07.11-.14.17-.21l.07-.07A5.11825,5.11825,0,0,1,11.85,5.78l.07-.07c.07-.06.14-.12.21-.17.11-.1.24-.19.36-.28a6.66426,6.66426,0,0,1,1.83-.91,5.43057,5.43057,0,0,1,.54-.15,2.87486,2.87486,0,0,1,.51-.1c.04-.01.08-.01.12-.02A2.9405,2.9405,0,0,1,16,4.03a4.18606,4.18606,0,0,1,1,0,2.9405,2.9405,0,0,1,.51.05c.04.01.08.01.12.02a2.87486,2.87486,0,0,1,.51.1,5.43057,5.43057,0,0,1,.54.15,6.66426,6.66426,0,0,1,1.83.91c.12.09.25.18.36.28.07.05.14.11.21.17a5.95007,5.95007,0,0,1,.71.71c.06.07.12.14.17.21.1.11.19.24.28.36a6.66426,6.66426,0,0,1,.91,1.83,5.43057,5.43057,0,0,1,.15.54,2.87925,2.87925,0,0,1,.1.51c.00995.04.00995.08.02.12a2.934,2.934,0,0,1,.05.51,4.17762,4.17762,0,0,1,0,1A2.93513,2.93513,0,0,1,23.42,12.01Z"></path>
                  <path fill="#b72732" d="M9,11a8.00539,8.00539,0,0,1,7.75-7.9873C16.66626,3.01007,16.58441,3,16.5,3a8,8,0,0,0,0,16c.08441,0,.16626-.01007.25-.0127A8.00544,8.00544,0,0,1,9,11Z"></path>
                </svg>

              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">Kiểm Soát Thời Gian</h3>
              <p className="text-gray-600 mb-4">Chọn thời gian phù hợp</p>
              <div className="space-y-2">
                <button className="w-full px-4 py-2 text-xl border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Bullet 
                </button>
                <button className="w-full px-4 py-2 text-xl border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Blitz 
                </button>
                <button className="w-full px-4 py-2 text-xl border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Rapid 
                </button>
              </div>
            </div>
          </div>

          {/* Learn */}
          <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200 hover:shadow-xl transition-shadow">
            <div className="text-center">
              <div className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Image
                  src="/assets/Brain.png"
                  width={64}
                  height={64}
                  alt="Icon"
                />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">Học Tập</h3>
              <p className="text-gray-600 mb-4">Cải thiện kỹ năng của bạn</p>
              <button className="w-full px-6 py-3 bg-black text-white rounded-lg hover:bg-blue-700 font-medium transition-colors">
                Câu đố
              </button>
              <button className="w-full mt-4 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors">
                Công cụ phân tích
              </button>
            </div>
          </div>
        </div>

        {/* Stats Section */}
        {user && (
          <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200 mb-12">
            <h3 className="text-xl font-semibold text-gray-900 mb-6">Thống Kê Của Bạn</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">1250</div>
                <div className="text-sm text-gray-600">Rapid Rating</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">45</div>
                <div className="text-sm text-gray-600">Trận thắng</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-red-600">23</div>
                <div className="text-sm text-gray-600">Trận thua</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-600">12</div>
                <div className="text-sm text-gray-600">Trận hòa</div>
              </div>
            </div>
          </div>
        )}

      </main>

    </div>
  );
}