import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import { DeploymentsPage, StatusPage } from './components';
import { Topbar } from './components/topbar';

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen">
        <Topbar />
        <Routes>
          <Route path="/" element={<StatusPage></StatusPage>} />
          <Route path="/deployments" element={<DeploymentsPage></DeploymentsPage>} />
          <Route path="/logs" element={<div> logs </div>} />
          <Route path="/config" element={<div> config </div>} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
