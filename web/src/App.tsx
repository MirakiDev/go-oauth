import { useEffect, useState } from "react";
import { loginWithGoogleProvider } from "./api/login";
import GridPattern from "./components/grid-pattern";
import { RainbowButton } from "./components/rainbow-button";
import { cn } from "./lib/utils";
import { FcGoogle } from "react-icons/fc";
import { logout } from "./api/logout";
import { getMe } from "./api/get-me";

type User = {
  id: string;
  email: string;
  picture: string;
};
export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userData, setUserData] = useState<User | null>(null);

  useEffect(() => {
    const saveUserData = async () => {
      const data = await getMe();
      if (data) {
        setUserData(data);
      }
    };

    const success = new URLSearchParams(window.location.search).get("success");

    if (success === "ok") {
      setIsLoggedIn(true);
      saveUserData();
    }
    
  }, []);

  const handleLogin = () => loginWithGoogleProvider();
  const handleLogout = () => logout();

  return (
    <div className="flex items-center justify-center h-screen overflow-hidden">
      <GridPattern
        width={30}
        height={30}
        x={-1}
        y={-1}
        strokeDasharray={"4 2"}
        className={cn(
          "[mask-image:radial-gradient(1400px_circle_at_center,white,transparent)]"
        )}
      />

      {userData && (
        <div className="flex flex-col gap-10 ">
          <div className="flex items-center justify-center gap-4">
            <p className="text-xl font-medium">Hi {userData.email} </p>
            {userData.picture && (
              <img
                src={userData.picture}
                alt="User Profile"
                className="rounded-full h-10 w-10 relative"
              />
            )}
          </div>
          <RainbowButton className="text-sm relative" onClick={handleLogout}>
            Log out
          </RainbowButton>
        </div>
      )}

      {!isLoggedIn && (
        <RainbowButton className="text-sm relative" onClick={handleLogin}>
          <FcGoogle size={20} className="mr-3" />
          Login with Google
        </RainbowButton>
      )}
    </div>
  );
}
