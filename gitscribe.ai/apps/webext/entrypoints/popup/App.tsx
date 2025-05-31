import { useState, useEffect } from 'react';
import { onAuthStateChanged, User } from 'firebase/auth';
import { auth } from '@/lib/firebase';
import { LoginForm } from './components/LoginForm';
import { Button } from '@workspace/ui/components/button';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import './App.css';

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const handleSignOut = async () => {
    try {
      await auth.signOut();
    } catch (error) {
      console.error('Error signing out:', error);
    }
  };

  if (loading) {
    return (
      <div className="w-80 h-96 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-2 text-sm text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="w-80 h-96 p-4">
        <LoginForm />
      </div>
    );
  }

  return (
    <div className="w-80 h-96 p-4">
      <Card>
        <CardHeader>
          <CardTitle>Welcome to GitScribe</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="text-sm">
            <p className="font-medium">Signed in as:</p>
            <p className="text-gray-600">{user.email}</p>
          </div>

          <div className="space-y-2">
            <Button className="w-full" variant="outline">
              Start Meeting Documentation
            </Button>
            <Button className="w-full" variant="outline">
              View Past Meetings
            </Button>
          </div>

          <Button
            onClick={handleSignOut}
            variant="ghost"
            className="w-full text-red-600 hover:text-red-700"
          >
            Sign Out
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}

export default App;
