import React, { useState, useEffect } from 'react';

import './App.css';

function App() {
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

    return (
        <div className="App">
          <header>
            <h1>Thunderdrone dashboard</h1>
          </header>
          <main>
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
          </main>
          <footer>

          </footer>
        </div>
    );
}

export default App;
