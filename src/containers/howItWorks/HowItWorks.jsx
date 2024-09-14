import { InstructionCard } from "../../components"
import { MdOutlineFileUpload } from "react-icons/md";
import { MdOutlineFileDownload } from "react-icons/md";
import { MdSearch } from "react-icons/md";

const HowItWorks = () => {
  return (
    <section>
      <div className="container flex flex-col items-center justify-center gap-8 bg-primary rounded p-6">
        <p className="text-4xl font-bold text-center text-white">How it works</p>
        <div className="flex flex-col gap-4 md:gap-8 md:grid md:grid-cols-3">
          <InstructionCard icon={<MdOutlineFileUpload size={40} className="text-primary rounded bg-white p-1.5 mb-3"/>} title="Upload" description="Upload a CSV or text file with your domain list."/>
          <InstructionCard  icon={<MdSearch size={40} className="text-primary rounded bg-white p-1.5 mb-3"/>} title="Search" description="Enter a keyword to search through the retrieved TXT records."/>
          <InstructionCard  icon={<MdOutlineFileDownload size={40} className="text-primary rounded bg-white p-1.5 mb-3"/>} title="Download" description="Export a file containing all domains with matching TXT records."/>
        </div>
      </div>
    </section>
  )
}

export default HowItWorks