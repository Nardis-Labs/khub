import React from 'react';
import { Root } from './app/Root/Root';

import { BrowserRouter as Router } from 'react-router-dom';
import { Theme } from '@carbon/react';
import { useSelector } from 'react-redux';
import { RootState } from './app/store';

function App() {
  const appTheme = useSelector((state: RootState) => state.appTheme);
  let carbonTheme: "white" | "g10" | "g90" | "g100" | undefined = undefined;
  if (appTheme.theme === 'white') {
    carbonTheme = 'white';
  } else if (appTheme.theme === 'dark') {
    carbonTheme = 'g100';
  }
  
  return (
      <Router>
        <div className={carbonTheme === 'g100' ? 'root-container-dark' : 'root-container'}>
          <Theme theme={carbonTheme} className={carbonTheme === 'g100' ? 'root-container-dark' : 'root-container'}>
            <Root/>
          </Theme>
        </div>
      </Router>
  );
}

export default App;
