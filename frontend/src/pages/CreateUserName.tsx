import { useState } from "react"
import { useParams, useNavigate } from "react-router";

const CreateUserName = () => {
    const [username, setUsername] = useState("");

    const navigator = useNavigate();
    const { id } = useParams();

    const setUser = async () => {
        await fetch("http://localhost:8080/createsession", {
            method: "POST",
            body: JSON.stringify({
                "username": username
            }),
            credentials: 'include' 
        })
        .then(() => {

            navigator(`/joinroom/${id}`)

        })
    }

    return <div>
            <div className="fieldWrapper">
                <input type="text" value={username} onChange={(e) => {setUsername(e.target.value)}} id="username" />
            </div>
            <div className="fieldWrapper">
                <button onClick={() => {setUser()}}>Create User</button>
            </div>
    </div>
}

export default CreateUserName