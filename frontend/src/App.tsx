import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { DeploymentsPage, EnvironementHealth, NavBar, StatusPage, Topbar } from './components';
import { ROUTES } from './lib';

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen pb-15 sm:pb-0">
        <Topbar>
          {/* Top navigation bar, on big screens */}
          <div className="flex-1 justify-around flex ml-10">
            <div className="flex-1 flex ">
              <NavBar className="hidden sm:flex bg-sidebar items-center" />
            </div>
            <EnvironementHealth></EnvironementHealth>
          </div>
        </Topbar>
        {/* Bottom navigation bar, on small screens */}
        <NavBar className="flex sm:hidden bg-sidebar h-14 border-t w-full fixed items-center justify-around bottom-0 left-0 right-0 z-50" />

        <Routes>
          <Route path="/" element={<Navigate to={ROUTES.DEPLOYMENTS}></Navigate>} />
          <Route path={ROUTES.DEPLOYMENTS} element={<DeploymentsPage></DeploymentsPage>} />
          <Route path={ROUTES.DEPLOYMENT(':id')} element={<DeploymentsPage></DeploymentsPage>} />
          <Route path={ROUTES.STATUS} element={<StatusPage></StatusPage>} />
          <Route path={ROUTES.LOGS} element={<div> logs </div>} />
          <Route path={ROUTES.CONFIG} element={<div> config </div>} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
