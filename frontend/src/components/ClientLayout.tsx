'use client';

import { useAuth } from '@/contexts/AuthContext';
import Navbar from '@/components/Navbar';

const ClientLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading, setUser } = useAuth();

  const handleLogout = async () => {
    try {
      await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/api/v1/auth/logout`, {
        method: 'POST',
        credentials: 'include'
      });
    } catch (error) {
      console.error('Logout error:', error);
    }
    
    setUser(null); 
    window.location.href = '/';
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-700 text-lg">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar user={user} onLogout={handleLogout} />
      <main className="pt-0">
        {children}
      </main>
    </div>
  );
};

export default ClientLayout;