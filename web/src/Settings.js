import React, {useState} from "react";
import Input from "./Input";
import Button from "./Button";
import useFetch from "./useFetch";

export default function Settings(props) {
    const [certificate, setCertificate] = useState()
    const [macaroon, setMacaroon] = useState()
    const [grpcUrl, setGrpcUrl] = useState()
    const {postFiles, loading} = useFetch("http://localhost:8080/api/")

    function handleFormSubmit(e) {
        e.preventDefault()

        const formData = new FormData();
        formData.append("implementation", "LND");
        formData.append("tlsFile", certificate);
        formData.append("macaroonFile", macaroon);
        formData.append("grpcAddress", grpcUrl);

        postFiles("nodesettings", formData)
            .then(r => console.log(r))
            .catch(e => console.log(e))
    }

    function handleCertificateChanged(e) {
        setCertificate(e.target.files[0])
    }

    function handleMacaroonChanged(e) {
        setMacaroon(e.target.files[0])
    }

    function handleGrpcUrlChanged(e) {
        setGrpcUrl(e.target.value)
    }

    return <>
        <form className="form" onSubmit={handleFormSubmit}>
            <p>
                Supply the necessary LND connection details
            </p>
            <Input placeholder="TLS Certificate" type="file" required onChange={handleCertificateChanged}/>
            <Input placeholder="Macaroon" type="file" required onChange={handleMacaroonChanged}/>
            <Input placeholder="GRPC Url (host:port)" type="text" required onChange={handleGrpcUrlChanged}/>
            <Button type="submit">Save</Button>
        </form>
    </>
}

