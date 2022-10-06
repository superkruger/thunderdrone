import React from 'react';
import {BrowserRouter, Routes, Route} from "react-router-dom";
import Navbar from "./Navbar";
import './App.css';
import Settings from "./Settings";

function App() {

    return (
        <BrowserRouter>
            <header>
                <nav className="navbar">
                    <Navbar/>
                </nav>
            </header>
            <div className="container">

              <main>
                  <Routes>
                      <Route exact path="/" element={<Dashboard/>}/>
                      <Route exact path="/settings" element={<Settings/>}/>
                  </Routes>
              </main>

              <footer>

              </footer>
            </div>
        </BrowserRouter>
    );
}

function Dashboard() {

    return <div className="dashboard-layout">
        <h1>Dashboard</h1>
    </div>
}

export default App;
