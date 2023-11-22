import {
  RouterProvider,
} from "react-router-dom";
import TurmsRouter from './TurmsRouter';



function App() {
  return (
    <div className="App">
      <h1>Turms</h1>
      <RouterProvider router={TurmsRouter} />
    </div>
  );
}

export default App;
