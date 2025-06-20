import { yupResolver } from "@hookform/resolvers/yup";
import Button from "../../components/button";
import Input from "../../components/input";

import { useForm, type SubmitHandler } from "react-hook-form";
import * as yup from "yup";
import { useMutation } from "@tanstack/react-query";

const HEADING = "Downtask";
const SUBHEADING = "Welcome back! Log-in to your account";
const FORGOT_PASSWORD = "Forgot Password?";
const SUBMIT_BTN = "Sign In";
const NEW_ACCOUNT = "Don't have an account?";
const SIGN_UP = "Sign Up";
const API = "http://localhost:3001/api";

interface IFormInput {
  email: string;
  password: string;
}

const schema = yup.object({
  email: yup.string().email().required(),
  password: yup.string().required().min(3).max(20),
});

const loginUser = async (data: IFormInput) => {
  const url = `${API}/v1/users/auth`;

  try {
    const res = await fetch(url, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (res.status === 400) {
      const body = await res.json();
      throw new Error(body.error || "invalid email or password");
    }

    if (res.status === 401) {
      const body = await res.json();
      throw new Error(body.error || "invalid email or password");
    }

    return res.json();
  } catch (error) {
    throw error;
  }
};

const SignIn = () => {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm({ resolver: yupResolver(schema) });

  const mutation = useMutation({
    mutationFn: loginUser,
    onSuccess: async (data) => {
      if (data.ok) {
        console.log("User logged in successfully"); // todo: remove this and navigate to login page
      }
    },
    onError: (resp) => {
      setError(
        "email",
        {
          message: resp.message,
          type: "validate",
        },
        { shouldFocus: true }
      );
    },
  });

  const formSubmitHandler: SubmitHandler<IFormInput> = async (data) => {
    mutation.mutate(data);
  };

  return (
    <div className="auth-container">
      <div className="heading-container">
        <h1 className="auth-heading">{HEADING}</h1>
        <h3 className="auth-subheading">{SUBHEADING}</h3>
      </div>

      <div className="form-container">
        <form onSubmit={handleSubmit(formSubmitHandler)}>
          <Input
            {...register("email", { required: true })}
            placeholder="username@downtask.com"
            className={`auth-input ${errors.email ? "error" : ""}`}
          />
          {errors.email ? (
            <span className="error-text">{errors.email.message}</span>
          ) : null}

          <div className="password-wrapper">
            <Input
              {...register("password", { required: true })}
              type="password"
              placeholder="your-strong-password"
              className={`auth-input ${errors.password ? "error" : ""}`}
            />
            {errors.password ? (
              <span className="error-text">{errors.password.message}</span>
            ) : null}
          </div>

          <p className="forgot-password">{FORGOT_PASSWORD}</p>

          <Button
            disabled={Object.keys(errors).length !== 0}
            type="submit"
            className="btn submit-btn"
          >
            {SUBMIT_BTN}
          </Button>
        </form>

        <p className="already-account-wrapper">
          <span>{NEW_ACCOUNT}</span>
          <span>{SIGN_UP}</span>
        </p>
      </div>
    </div>
  );
};

export default SignIn;
