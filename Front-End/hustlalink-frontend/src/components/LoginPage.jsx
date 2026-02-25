// src/pages/SignupPage.jsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";

function SignupPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSignup = () => {
    navigate("/main");
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>HustlaLink Sign Up</h1>
      <input 
        type="text" 
        placeholder="Full Name" 
        value={name}
        onChange={(e) => setName(e.target.value)}
      /><br/><br/>
      <input 
        type="email" 
        placeholder="Email" 
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      /><br/><br/>
      <input 
        type="password" 
        placeholder="Password" 
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      /><br/><br/>
      <button onClick={handleSignup}>Sign Up</button>
    </div>
  );
}

export default SignupPage;