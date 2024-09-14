import { Header, FileUpload, FileDownload, HowItWorks, Footer } from './containers'
import { NavBar, CTA, InstructionCard, FileInput } from './components'

const App = () => {
  return (
    <div className="app">
      <NavBar />
      <Header />
      <FileUpload />
      <FileDownload />
      <HowItWorks />
      <Footer />
    </div>
  )
}

export default App