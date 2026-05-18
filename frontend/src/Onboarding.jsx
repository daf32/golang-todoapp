import React from 'react';

const Onboarding = () => {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      height: '100vh',
      backgroundColor: '#f0f2f5',
      fontFamily: 'Arial, sans-serif',
      textAlign: 'center',
      padding: '20px'
    }}>
      <div style={{
        width: '100%',
        maxWidth: '400px',
        backgroundColor: 'white',
        borderRadius: '20px',
        padding: '40px',
        boxShadow: '0 4px 8px rgba(0,0,0,0.1)'
      }}>
        <div style={{ marginBottom: '40px' }}>
          <svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="circleGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style={{stopColor: '#8e44ad', stopOpacity: 1}} />
                <stop offset="100%" style={{stopColor: '#3498db', stopOpacity: 1}} />
              </linearGradient>
            </defs>
            <circle cx="50" cy="50" r="40" fill="url(#circleGradient)" />
            <path d="M 30 50 Q 50 30 70 50" stroke="white" strokeWidth="4" fill="none" />
            <path d="M 30 60 Q 50 40 70 60" stroke="white" strokeWidth="4" fill="none" />
          </svg>
        </div>
        <h1 style={{
          fontSize: '28px',
          fontWeight: 'bold',
          color: '#333',
          margin: '0 0 10px'
        }}>
          Manage your daily tasks
        </h1>
        <p style={{
          fontSize: '16px',
          color: '#666',
          marginBottom: '40px'
        }}>
          Team and Project management with solution providing App
        </p>
        <button style={{
          backgroundColor: '#8e44ad',
          color: 'white',
          border: 'none',
          borderRadius: '10px',
          padding: '15px 30px',
          fontSize: '16px',
          fontWeight: 'bold',
          cursor: 'pointer',
          boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
        }}>
          Get Started
        </button>
      </div>
    </div>
  );
};

export default Onboarding;
