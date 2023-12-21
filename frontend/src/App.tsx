import {
  RouterProvider,
} from "react-router-dom";
import TurmsRouter from './TurmsRouter';



function App() {
  return (
    <div className="App">
      <div className="navbar">
        <a className="navbar-head" href="/">Turms</a> 
        <a href="/session">Session</a> 
        <a href="/join">Join</a> 
        <a href="/about">About</a>
      </div>
      <hr />
      <RouterProvider router={TurmsRouter} />
    </div>
  );
}

export default App;
