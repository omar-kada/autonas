import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { DeploymentsPage, StatusPage } from './components';
import { Topbar } from './components/topbar';
import { ROUTES } from './lib';

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen">
        <Topbar />
        <Routes>
          <Route path="/" element={<StatusPage></StatusPage>} />
          <Route path={ROUTES.STATUS} element={<StatusPage></StatusPage>} />
          <Route path={ROUTES.DEPLOYMENTS} element={<DeploymentsPage></DeploymentsPage>} />
          <Route path={ROUTES.DEPLOYMENT(':id')} element={<DeploymentsPage></DeploymentsPage>} />
          <Route path={ROUTES.LOGS} element={<div> logs </div>} />
          <Route path={ROUTES.CONFIG} element={<div> config </div>} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
