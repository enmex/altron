import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { SignIn } from './pages/Auth/SignIn';
import { Home } from './pages/Home/Home';
import { Service } from './pages/Service/Service';
import { useAppSelector } from './app/store/hooks';
import { CART_PATH, INDEX_PATH, LOGS_PATH, PCAP_PATH, SERVICE_PATH, SESSION_PATH, SIGN_IN_PATH, WORKSPACE_PATH } from './config/constants';
import { Workspace } from './pages/Workspace/Workspace';
import { Session } from './pages/Session/Session';
import { ProtectedLayout } from './layouts/ProtectedLayout';
import { Logs } from './pages/Logs/Logs';
import { NotificationToaster } from './components/organisms/NotificationToaster';
import { Cart } from './pages/Cart/Cart';
import { Pcap } from './pages/Pcap/Pcap';
import { ServiceFallbackLayout } from './layouts/ServiceFallbackLayout';
import { AuthLayout } from './layouts/AuthLayout';
import { NotFound } from './pages/404/404';
import { PcapWorkspaceFallbackLayout } from './layouts/PcapWorkspaceFallbackLayout';
import { WorkspaceFallbackLayout } from './layouts/WorkspaceFallbackLayout';

const BaseRouter = () => {
  return (
    <>
    <NotificationToaster />
    <BrowserRouter>
        <Routes>
          <Route path={SIGN_IN_PATH} element={
            <AuthLayout>
              <SignIn />
            </AuthLayout>
          } />
          <Route path={INDEX_PATH} element={
            <ProtectedLayout>
              <Home/>
            </ProtectedLayout>
          } />
          <Route path={SERVICE_PATH} element={
            <ProtectedLayout>
              <ServiceFallbackLayout>
                <Service />
              </ServiceFallbackLayout>
            </ProtectedLayout>
          } />
          <Route path={WORKSPACE_PATH} element={
            <ProtectedLayout>
              <WorkspaceFallbackLayout>
                <Workspace />
              </WorkspaceFallbackLayout>
            </ProtectedLayout>
          } />
          <Route path={SESSION_PATH} element={
            <ProtectedLayout>
              <Session />
            </ProtectedLayout>
          } />
          <Route path={LOGS_PATH} element={
            <ProtectedLayout>
              <ServiceFallbackLayout>
                <Logs />
              </ServiceFallbackLayout>
            </ProtectedLayout>
          } />
          <Route path={CART_PATH} element={
            <ProtectedLayout>
              <WorkspaceFallbackLayout>
                <Cart />
              </WorkspaceFallbackLayout>
            </ProtectedLayout>
          } />
          <Route path={PCAP_PATH} element={
            <ProtectedLayout>
              <PcapWorkspaceFallbackLayout>
                <Pcap />
              </PcapWorkspaceFallbackLayout>
            </ProtectedLayout>
          }
          />
          <Route path={"/"} element={ <Navigate to={"/home"}/> } />
          <Route path="*" element={<NotFound />}/>
        </Routes>
    </BrowserRouter>
    </>
  );
}

const App = () => {
  const theme = useAppSelector(state => state.rootReducer.theme);

  return (
    <div 
      className='fixed w-screen h-screen'
      style={{
        backgroundColor: theme.primary
      }}
    >
      <BaseRouter/>
    </div>
  );
};

export default App;
