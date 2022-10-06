import React, {useState, useEffect} from "react";
import Button from "./Button";
import useFetch from "./useFetch";
import NodeSettings from "./NodeSettings";

export default function Settings(props) {
    const [localNodes, setLocalNodes] = useState([])

    const {get, loading} = useFetch("http://localhost:8080/api/")

    useEffect(() => {
        get("nodesettings")
            .then(d => setLocalNodes(d))
            .catch(e => console.log(e))
    }, [])

    function handleAddClicked() {
        setLocalNodes(prev => [...prev, {}])
    }

    return <>
        <div className="settings-layout">
            <h1>Local Node Settings</h1>
            <p>Supply the necessary LND connection details</p>
            <div className="node-settings-grid">
                {
                    localNodes.map(localNode => {
                        return (
                            <NodeSettings localNode={localNode}/>
                        )
                    })
                }
            </div>
        </div>
        <div>
            <Button onClick={handleAddClicked}>Add Local Node</Button>
        </div>
    </>
}

