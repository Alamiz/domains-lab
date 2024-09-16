import 'react-toastify/dist/ReactToastify.css';
import { Header, FileUpload, FileDownload, HowItWorks, Footer } from './containers'
import { NavBar } from './components'
import { ToastContainer } from 'react-toastify';

const App = () => {

  return (
    <div className="app">
      <NavBar />
      <Header />
      <FileUpload />
      <FileDownload />
      <HowItWorks />
      <Footer />
      <ToastContainer />
    </div>
  )
}

export default App