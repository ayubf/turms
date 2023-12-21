import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

const Session = () => {
  function getRandomInt() {
    return Math.floor(Math.random() * 99999);
  }

  const navigator = useNavigate();

  const [sessionExists, setSessionExist] = useState(false);
  const [username, setUsername] = useState("");
  const [inputUsername, setInputUsername] = useState("");
  const [sessExpire, setSessExpire] = useState<Date>();

  const checkSession = async () => {
    try {
      const response = await fetch(`http://localhost:8080/getsession`, {
        credentials: "include",
      });
      const data = await response.json();

      if (data.hasOwnProperty("iss")) {
        setSessionExist(true);
        setUsername(data["iss"]);
        setSessExpire(new Date(parseInt(data["exp"]) * 1000));
      } else {
        setUsername("Guest" + getRandomInt());
      }
    } catch (error) {
      console.error("Error fetching session:", error);
    }
  };

  const createSession = async () => {
    await fetch(`http://localhost:8080/createsession`, {
      method: "POST",
      body: JSON.stringify({
        username: inputUsername || username,
      }),
      credentials: "include",
    }).then(() => {
        navigator(0);
    })
  };

  useEffect(() => {
    checkSession();
    setInputUsername(""); // Clear inputUsername when fetching the session
  }, [username]);

  return (
    <div>
      <h1>Session</h1>
      <div>
        {sessionExists ? (
          <div>
            <p>Current Session As: {username}</p>
            <p>Session Expires At: {sessExpire?.toLocaleString()} </p>
          </div>
        ) : (
          <></>
        )}
      </div>
      <div className="fieldWrapper">
        <input
          type="text"
          value={inputUsername}
          onChange={(e) => setInputUsername(e.target.value)}
          id="username"
        />
      </div>
      <div className="fieldWrapper">
        <button onClick={createSession}>Create New Session</button>
      </div>
    </div>
  );
};

export default Session;
