'use client';

import React from 'react';
import Link from 'next/link';
import { Mail, User, Lock, Eye, EyeOff, ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation'; 

interface RegisterData {
  username: string;
  email: string;
  password: string;
  skillLevel: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

const RegisterPage: React.FC = () => {
  const router = useRouter();
  const [step, setStep] = React.useState<'choice' | 'form'>('choice');
  const [formData, setFormData] = React.useState<RegisterData>({
    username: '',
    email: '',
    password: '',
    skillLevel: 'beginner'
  });
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState('');
  const [showPassword, setShowPassword] = React.useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    if (error) setError('');
  };

  const getEloFromSkillLevel = (level: string) => {
    switch (level) {
      case 'beginner': return 800;
      case 'intermediate': return 1200;
      case 'advanced': return 1600;
      case 'expert': return 2000;
      default: return 800;
    }
  };

  const validateForm = () => {
    if (!formData.username.trim()) {
      setError('Username is required');
      return false;
    }
    if (formData.username.length < 3) {
      setError('Username must be at least 3 characters');
      return false;
    }
    if (!formData.email.trim()) {
      setError('Email is required');
      return false;
    }
    if (!/\S+@\S+\.\S+/.test(formData.email)) {
      setError('Please enter a valid email address');
      return false;
    }
    if (!formData.password) {
      setError('Password is required');
      return false;
    }
    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters');
      return false;
    }
    return true;
  };

 const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();

  if (!validateForm()) return;

  setLoading(true);
  setError('');

  try {
    const registerData = {
      username: formData.username,
      email: formData.email,
      password: formData.password,
      elo: getEloFromSkillLevel(formData.skillLevel)
    };

    console.log('Register data to send:', registerData);

    const response = await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/api/v1/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(registerData)
    });

    const result = await response.json();

    if (response.ok) {
      alert('Registration successful!');
      router.replace('/login');
      setFormData({
        username: '',
        email: '',
        password: '',
        skillLevel: 'beginner'
      });
      setStep('choice'); 
    } else {
      setError(result.error || 'Registration failed');
    }

  } catch (error) {
    console.error('Registration error:', error);
    setError('Registration failed. Please try again.');
  } finally {
    setLoading(false);
  }
};

  const handleGoogleRegister = () => {
    alert('Google registration will be implemented later!');
  };

  const handleContinueWithEmail = () => {
    setStep('form');
    setError('');
  };

  // Step 1: Choice between Google and Email
  if (step === 'choice') {
    return (
      <div 
        className="min-h-screen bg-cover bg-center bg-no-repeat flex items-center justify-center px-4"
        style={{ backgroundImage: 'url(/register-background.png)' }}
      >
        <div className="w-full max-w-md">
          <div className="bg-white/95 backdrop-blur-sm rounded-2xl p-8 shadow-2xl">
            
            <Link href="/login">
              <button className="mb-6 flex items-center text-gray-600 hover:text-gray-800 transition-colors">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Login
              </button>
            </Link>

            <div className="text-center mb-8">
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Kangyoo Chess</h1>
          
            </div>

            <div className="space-y-4">
              <button
                onClick={handleGoogleRegister}
                disabled={loading}
                className="w-full py-4 bg-white border-2 border-gray-300 hover:border-gray-400 hover:shadow-lg text-gray-800 font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-gray-400 disabled:opacity-50 disabled:transform-none flex items-center justify-center"
              >
                <svg className="w-5 h-5 mr-3" viewBox="0 0 24 24">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                Continue with Google
              </button>

              <div className="flex items-center">
                <div className="flex-1 border-t border-gray-300"></div>
                <span className="px-4 text-sm text-gray-500">or</span>
                <div className="flex-1 border-t border-gray-300"></div>
              </div>

              <button
                onClick={handleContinueWithEmail}
                disabled={loading}
                className="w-full py-4 bg-black hover:bg-gray-800 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-gray-400 disabled:opacity-50 disabled:transform-none flex items-center justify-center"
              >
                <Mail className="h-5 w-5 mr-3" />
                Continue with Email
              </button>
            </div>

            <div className="mt-8 text-center">
              <p className="text-gray-600 text-sm">
                Already have an account?{' '}
                <Link href="/login" className="text-blue-500 hover:text-blue-600 font-semibold transition-colors">
                  Sign in
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Step 2: Email Registration Form
  return (
    <div 
      className="min-h-screen bg-cover bg-center bg-no-repeat flex items-center justify-center px-4"
      style={{ backgroundImage: 'url(/register-background.png)' }}
    >
      <div className="w-full max-w-md">
        <div className="bg-white/95 backdrop-blur-sm rounded-2xl p-8 shadow-2xl">

          <button
            onClick={() => setStep('choice')}
            className="mb-6 flex items-center text-gray-600 hover:text-gray-800 transition-colors"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </button>

          <div className="text-center mb-8">

            <h1 className="text-3xl font-bold text-gray-900 mb-2">Create Account</h1>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl">
              <p className="text-red-600 text-sm text-center">{error}</p>
            </div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">

            {/* Username */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <User className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                name="username"
                placeholder="Username"
                value={formData.username}
                onChange={handleInputChange}
                required
                disabled={loading}
                className="w-full pl-12 pr-4 py-4 bg-gray-50 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent transition-all duration-300 disabled:opacity-50"
              />
            </div>

            {/* Email */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <Mail className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="email"
                name="email"
                placeholder="Email address"
                value={formData.email}
                onChange={handleInputChange}
                required
                disabled={loading}
                className="w-full pl-12 pr-4 py-4 bg-gray-50 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent transition-all duration-300 disabled:opacity-50"
              />
            </div>

            {/* Password */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <Lock className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type={showPassword ? "text" : "password"}
                name="password"
                placeholder="Password"
                value={formData.password}
                onChange={handleInputChange}
                required
                disabled={loading}
                className="w-full pl-12 pr-12 py-4 bg-gray-50 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent transition-all duration-300 disabled:opacity-50"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                disabled={loading}
                className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-gray-600 transition-colors"
              >
                {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
              </button>
            </div>

            {/* Skill Level */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Choose your skill level
              </label>
              <div className="grid grid-cols-2 gap-3">
                {[
                  { value: 'beginner', label: 'Beginner', elo: '800', icon: '♟' },
                  { value: 'intermediate', label: 'Intermediate', elo: '1200', icon: '♞' },
                  { value: 'advanced', label: 'Advanced', elo: '1600', icon: '♜' },
                  { value: 'expert', label: 'Expert', elo: '2000', icon: '♛' }
                ].map((level) => (
                  <button
                    key={level.value}
                    type="button"
                    onClick={() => setFormData({ ...formData, skillLevel: level.value as 'beginner' | 'intermediate' | 'advanced' | 'expert'})}
                    disabled={loading}
                    className={`p-4 rounded-xl border-2 transition-all text-left disabled:opacity-50 ${
                      formData.skillLevel === level.value
                        ? 'border-black bg-gray-100'
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                  >
                    <div className="text-lg mb-1 text-gray-900">{level.icon} {level.label}</div>
                    <div className="text-xs text-gray-500">ELO: {level.elo}</div>
                  </button>
                ))}
              </div>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-4 bg-black hover:bg-gray-800 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-gray-400 disabled:opacity-50 disabled:transform-none"
            >
              {loading ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                  Creating Account...
                </div>
              ) : (
                'Create Account'
              )}
            </button>

          </form>

          {/* Login Link */}
          <div className="mt-6 text-center">
            <p className="text-gray-600 text-sm">
              Already have an account?{' '}
              <Link href="/login" className="text-blue-500 hover:text-blue-600 font-semibold transition-colors">
                Sign in
              </Link>
            </p>
          </div>

        </div>
      </div>
    </div>
  );
};

export default RegisterPage;