import React from 'react';

const Dashboard = () => {
  return (
    <div style={{
      backgroundColor: '#f0f2f5',
      fontFamily: 'Arial, sans-serif',
      padding: '20px',
      minHeight: '100vh'
    }}>
      <div style={{
        maxWidth: '400px',
        margin: '0 auto',
        backgroundColor: '#8e44ad',
        color: 'white',
        borderRadius: '20px',
        padding: '20px',
        boxShadow: '0 4px 8px rgba(0,0,0,0.1)'
      }}>
        <header style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '20px'
        }}>
          <div>
            <h1 style={{ margin: 0, fontSize: '24px' }}>Hi Ghulam</h1>
            <p style={{ margin: 0, fontSize: '14px', opacity: 0.8 }}>6 Tasks are pending</p>
          </div>
          <img
            src="https://i.pravatar.cc/50"
            alt="Profile"
            style={{ borderRadius: '50%', width: '50px', height: '50px' }}
          />
        </header>

        <div style={{
          backgroundColor: 'rgba(255, 255, 255, 0.2)',
          borderRadius: '15px',
          padding: '15px',
          marginBottom: '20px'
        }}>
          <h2 style={{ margin: '0 0 5px', fontSize: '18px' }}>Mobile App Design</h2>
          <p style={{ margin: '0 0 10px', fontSize: '12px', opacity: 0.8 }}>Mike and Ann</p>
          <button style={{
            backgroundColor: 'white',
            color: '#8e44ad',
            border: 'none',
            borderRadius: '10px',
            padding: '10px 20px',
            fontWeight: 'bold'
          }}>View</button>
        </div>

        <div style={{
          backgroundColor: 'rgba(255, 255, 255, 0.2)',
          borderRadius: '15px',
          padding: '15px',
          marginBottom: '20px'
        }}>
          <h2 style={{ margin: '0 0 10px', fontSize: '18px' }}>Monthly Review</h2>
          <button style={{
            backgroundColor: 'white',
            color: '#8e44ad',
            border: 'none',
            borderRadius: '10px',
            padding: '10px 20px',
            fontWeight: 'bold'
          }}>View</button>
        </div>

        <div style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr',
          gap: '15px',
          textAlign: 'center'
        }}>
          <div style={{ backgroundColor: 'rgba(255, 255, 255, 0.2)', borderRadius: '15px', padding: '20px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold' }}>22</div>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>Done</div>
          </div>
          <div style={{ backgroundColor: 'rgba(255, 255, 255, 0.2)', borderRadius: '15px', padding: '20px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold' }}>7</div>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>In progress</div>
          </div>
          <div style={{ backgroundColor: 'rgba(255, 255, 255, 0.2)', borderRadius: '15px', padding: '20px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold' }}>10</div>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>Ongoing</div>
          </div>
          <div style={{ backgroundColor: 'rgba(255, 255, 255, 0.2)', borderRadius: '15px', padding: '20px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold' }}>12</div>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>Waiting for review</div>
          </div>
        </div>

        <nav style={{
          display: 'flex',
          justifyContent: 'space-around',
          marginTop: '20px',
          padding: '10px 0',
          backgroundColor: 'rgba(255, 255, 255, 0.2)',
          borderRadius: '15px'
        }}>
          <div style={{fontSize: '24px'}}>🏠</div>
          <div style={{fontSize: '24px'}}>📁</div>
          <div style={{fontSize: '24px'}}>👤</div>
          <div style={{fontSize: '24px'}}>⚙️</div>
        </nav>
      </div>
    </div>
  );
};

export default Dashboard;
