import Button from "../../components/button";
import Input from "../../components/input";
import * as yup from "yup";
import { type SubmitHandler, useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import { useMutation } from "@tanstack/react-query";

interface IFormInput {
  email: string;
  password: string;
}

const HEADING = "Downtask";
const SUBHEADING = "Create a new account";
const FORGOT_PASSWORD = "Forgot Password?";
const ALREADY_ACCOUNT = "Already have an account?";
const API_URL = "http://localhost:3001/api";

const schema = yup.object({
  email: yup.string().email().required(),
  password: yup.string().required().min(3).max(20),
});

const Signup = () => {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm({ resolver: yupResolver(schema) });

  const sendRegister = async (data: IFormInput) => {
    const url = `${API_URL}/v1/users/register`;

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
      return res.json();
    } catch (error) {
      throw error;
    }
  };

  const mutation = useMutation({
    mutationFn: sendRegister,
    onSuccess: async (data) => {
      if (data.ok) {
        console.log("User account created");
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

  const onSubmit: SubmitHandler<IFormInput> = async (data) => {
    mutation.mutate(data);
  };

  return (
    <div className="auth-container">
      <div className="heading-container">
        <h1 className="auth-heading">{HEADING}</h1>
        <h3 className="auth-subheading">{SUBHEADING}</h3>
      </div>

      <div className="form-container">
        <form onSubmit={handleSubmit(onSubmit)}>
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
            type="submit"
            disabled={Object.keys(errors).length !== 0}
            className="btn submit-btn"
          >
            Sign Up
          </Button>
        </form>
        <p className="already-account-wrapper">
          <span>{ALREADY_ACCOUNT}</span>
          <span>Sign In</span>
        </p>
      </div>
    </div>
  );
};

export default Signup;
