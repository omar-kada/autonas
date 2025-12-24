import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { DeploymentsPage, NavBar, StatusPage } from './components';
import { EnvironementStats } from './components/environement-stats';
import { Topbar } from './components/topbar';
import { Separator } from './components/ui/separator';
import { ROUTES } from './lib';

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen pb-15 sm:pb-0">
        <Topbar>
          {/* Top navigation bar, on big screens */}
          <NavBar className="hidden sm:flex bg-sidebar h-12 my-1 items-center" />
        </Topbar>
        {/* Bottom navigation bar, on small screens */}
        <NavBar className="flex sm:hidden bg-sidebar py-2 h-14 border-t w-full fixed items-center justify-around bottom-0 left-0 right-0 z-50" />

        <EnvironementStats className="flex"></EnvironementStats>
        <Separator orientation="horizontal"></Separator>

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
