import {
  RouterProvider,
} from "react-router-dom";
import TurmsRouter from './TurmsRouter';



function App() {
  return (
    <div className="App">
      <div className="navbar">
        <a className="navbar-head" href="/">Turms</a> <p>/</p>
        <a href="/session">Session</a> <p>/</p>
        <a href="/about">About</a>
      </div>
      <hr />
      <RouterProvider router={TurmsRouter} />
    </div>
  );
}

export default App;
