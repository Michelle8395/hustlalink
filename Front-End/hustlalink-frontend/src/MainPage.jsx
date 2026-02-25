// src/pages/MainPage.jsx
import { useState } from "react";
import JobList from "./components/JobList";
import UserTypeSelect from "./components/UserTypeSelect";

const mockJobs = [
  { id: 1, title: "Laundry Service", type: "Laundry", name: "Aisha", rating: 5, reviews: 10 },
  { id: 2, title: "Carpentry Work", type: "Carpentry", name: "Peter", rating: 4, reviews: 7 },
  { id: 3, title: "Plumbing Fix", type: "Plumber", name: "John", rating: 3, reviews: 5 },
  { id: 4, title: "Gardening Service", type: "Gardener", name: "Mary", rating: 5, reviews: 12 },
];

function MainPage() {
  const [location, setLocation] = useState("");
  const [userType, setUserType] = useState("job-seeker");

  return (
    <div style={{ padding: 20 }}>
      <h1>HustlaLink Main</h1>
      <input 
        type="text" 
        placeholder="Enter your location" 
        value={location}
        onChange={(e) => setLocation(e.target.value)}
      /><br/><br/>
      
      <UserTypeSelect userType={userType} setUserType={setUserType} /><br/><br/>

      <JobList jobs={mockJobs} />
    </div>
  );
}

export default MainPage;