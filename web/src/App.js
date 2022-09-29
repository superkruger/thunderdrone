import React, { useState, useEffect } from 'react';
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
    const [nodes, setNodes] = useState([]);

    useEffect(() => {
        console.log("fetching nodes")
        fetch('http://localhost:8080/nodes')
            .then((response) => response.json())
            .then((data) => {
                console.log(data);
                setNodes(data);
            })
            .catch((err) => {
                console.log(err.message);
            });
    }, []);

    return <>
        {
            nodes.map(node => {
                return (
                    <div key={node.id}>
                        <h2>ID: {node.id}</h2>
                        <p>{node.address}:{node.port}</p>
                    </div>
                )
            })
        }
    </>
}

export default App;
