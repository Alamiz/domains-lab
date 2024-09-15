import { Header, FileUpload, FileDownload, HowItWorks, Footer } from './containers'
import { NavBar, CTA } from './components'
import { useState } from 'react'

const App = () => {
  const [isFileProcessed, setIsFileProcessed] = useState(false)

  return (
    <div className="app">
      <NavBar />
      <Header />
      <FileUpload setIsFileProcessed={setIsFileProcessed}/>
      <FileDownload isFileProcessed={isFileProcessed}/>
      <HowItWorks />
      <Footer />
    </div>
  )
}

export default App