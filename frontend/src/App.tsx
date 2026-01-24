import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import {
  ConfigPage,
  DeploymentsPage,
  EnvironementHealth,
  NavBar,
  StatusPage,
  Topbar,
} from './components';
import { ROUTES } from './lib';

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen pb-12 md:pb-0">
        <Topbar>
          {/* Top navigation bar, on big screens */}
          <div className="flex">
            <NavBar className="hidden md:flex bg-sidebar items-center flex-1" />
          </div>
          <EnvironementHealth></EnvironementHealth>
          {/* <div className="flex justify-before flex-row-reverse ml-10">
          </div> */}
        </Topbar>
        {/* Bottom navigation bar, on small screens */}
        <NavBar className="flex md:hidden bg-sidebar h-12 border-t w-full fixed items-center justify-around bottom-0 left-0 right-0 z-50" />

        <Routes>
          <Route path="/" element={<Navigate to={ROUTES.DEPLOYMENTS}></Navigate>} />
          <Route path={ROUTES.DEPLOYMENTS} element={<DeploymentsPage />} />
          <Route path={ROUTES.DEPLOYMENT(':id')} element={<DeploymentsPage />} />
          <Route path={ROUTES.STATUS} element={<StatusPage />} />
          <Route path={ROUTES.LOGS} element={<div> logs </div>} />
          <Route path={ROUTES.CONFIG} element={<ConfigPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
