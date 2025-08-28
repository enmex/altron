import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { store } from './app/store/store';
import { Provider } from 'react-redux';
import { I18nextProvider } from 'react-i18next';
import "flag-icons/css/flag-icons.min.css"
import i18n from './i18';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <Provider store={store}>
    <I18nextProvider i18n={i18n}>
      <App />
    </I18nextProvider>
  </Provider>
);