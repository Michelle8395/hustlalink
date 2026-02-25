function UserTypeSelect({ userType, setUserType }) {
  return (
    <div>
      <label>
        <input 
          type="radio" 
          value="job-seeker" 
          checked={userType === "job-seeker"} 
          onChange={(e) => setUserType(e.target.value)}
        /> Job Seeker
      </label>
      {" "}
      <label>
        <input 
          type="radio" 
          value="employer" 
          checked={userType === "employer"} 
          onChange={(e) => setUserType(e.target.value)}
        /> Employer
      </label>
    </div>
  );
}

export default UserTypeSelect;