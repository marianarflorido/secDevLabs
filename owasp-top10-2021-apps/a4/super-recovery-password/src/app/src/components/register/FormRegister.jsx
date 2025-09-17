import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

import FormTitle from "../titles/FormTitle";
import FormInput from "../inputs/FormInput";
import FormButton from "../buttons/FormButton";

import BoxError from "../error/BoxError";

import { RegisterService } from "../../services/requests";

const FormRegister = () => {

  const navigate = useNavigate();
  const [state, setState] = useState("register");
  const [login, setLogin] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");
  const [registerStatus, setRegisterStatus] = useState(true);
  const [message, setMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    setMessage("");
    setRegisterStatus(true);

    if (!login || !email || !password || !repeatPassword) {
      setMessage("All fields are required.");
      setRegisterStatus(false);
      return;
    }

    if (password !== repeatPassword) {
      setMessage("Passwords do not match");
      setRegisterStatus(false);
      return;
    }

    try {
      const register = await RegisterService({
        login,
        password,
        email,
      }); 
      if (register.message === "success") {
          navigate("/login", { msg: "User registered successfully!" });
      } else {
          setMessage(
            "User or email already exists or a problem has occurred. Try again!"
          );
          setRegisterStatus(false);
          setState("register");
      }
    } catch (err) {}
  };
      
  return (
    <form method="POST" onSubmit={handleSubmit}>
      <FormTitle title="Register" />
      {message && <BoxError message={message} />}
      
      <FormInput 
        placeholder="login" 
        type="text" 
        value={login} 
        setValue={setLogin} 
      />
      <FormInput 
        placeholder="email" 
        type="email" 
        value={email}
        setValue={setEmail} 
      />
      <FormInput
        placeholder="password"
        type="password"
        value={password}
        setValue={setPassword}
      />
      <FormInput
        placeholder="repeat password"
        type="password"
        value={repeatPassword}
        setValue={setRepeatPassword}
      />
      <FormButton type="submit" text="Register" />
    </form>
  );
};

export default FormRegister;
