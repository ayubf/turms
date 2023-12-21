import React, {useEffect, useState} from 'react';
import { useParams, redirect } from 'react-router';

const Room = () => {
    const {id} = useParams();
    const [sessionExists, setSessionExist] = useState(false);
    const [username, setUsername] = useState("");

    const checkSession = async () => {
        await fetch("http://localhost:8080/getsession", {
            credentials: 'include' 
        })
        .then(i => i.json())
        .then((data) => {
            if (data.hasOwnProperty("iss")) {
                setSessionExist(true); 
                setUsername(data["iss"])
            } else {
                redirect("/")
            }
        })
    }

    useEffect(() => {
        checkSession()
        console.log(sessionExists)
    })

    return <div>
        <h1>Room</h1>
        <p>{id}</p>
        <p>{username}</p>
    </div>
}

export default Room;