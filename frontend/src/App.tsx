import { useEffect, useState } from 'react';
import { BrowserRouter, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import {
  ConfigPage,
  DeploymentsPage,
  EnvironementHealth,
  ErrorAlert,
  LoginPage,
  NavBar,
  RegisterPage,
  StatusPage,
  Topbar,
} from './components';
import { useRegistered, useUser } from './hooks';
import { ROUTES } from './lib';

function RouteBasedTopBar() {
  const { data: isRegistered, isPending, error } = useRegistered();
  const { data: user, isPending: userPending, error: userError } = useUser();
  const navigate = useNavigate();
  const location = useLocation();

  const [showTopBar, setShowTopBar] = useState(!!user);

  useEffect(() => {
    setShowTopBar(!!isRegistered && !!user);
  }, [setShowTopBar, isRegistered, user]);

  useEffect(() => {
    if (isPending) {
      return;
    } else if (!isRegistered) {
      navigate(ROUTES.REGISTER);
    } else {
      if ((location.pathname = ROUTES.REGISTER)) {
        navigate(ROUTES.DEPLOYMENTS);
      }
      if (userPending) {
        return;
      } else if (!user) {
        navigate(ROUTES.LOGIN);
      } else {
        if ((location.pathname = ROUTES.LOGIN)) {
          navigate(ROUTES.DEPLOYMENTS);
        }
      }
    }
  }, [isRegistered, user, isPending, userPending]);

  const mergedError = error ?? userError;
  return (
    <>
      {showTopBar && (
        <>
          <Topbar>
            {/* Top navigation bar, on big screens */}
            <div className="flex">
              <NavBar className="hidden md:flex bg-sidebar items-center flex-1" />
            </div>
            <EnvironementHealth></EnvironementHealth>
          </Topbar>
          {/* Bottom navigation bar, on small screens */}
          <NavBar className="flex md:hidden bg-sidebar h-12 border-t w-full fixed items-center justify-around bottom-0 left-0 right-0 z-50" />
        </>
      )}
      <ErrorAlert title={mergedError?.message ?? null} />
    </>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-col h-screen pb-12 md:pb-0">
        <RouteBasedTopBar />
        <Routes>
          <Route path={ROUTES.ROOT} element={<Navigate to={ROUTES.DEPLOYMENTS}></Navigate>} />
          <Route path={ROUTES.REGISTER} element={<RegisterPage />} />
          <Route path={ROUTES.LOGIN} element={<LoginPage />} />
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
