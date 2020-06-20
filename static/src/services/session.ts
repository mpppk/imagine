import jwt_decode from 'jwt-decode';
import { JWTClaim } from '../models/models';
import { sleep } from '../util';

const getPasswordHash = (p: string) => p; // FIXME
export const requestSignIn = async (_email: string, password: string) => {
  const jwt = generateDummyJwt();
  const claim = jwt_decode<JWTClaim>(jwt);

  getPasswordHash(password);
  await sleep(1000);
  return {
    jwt,
    user: claim,
  };
};

const generateDummyJwt = () => {
  const header = 'ew0KICAidHlwIjogIkpXVCIsDQogICJhbGciOiAiSFMyNTYiDQp9';
  const payload =
    'ew0KICAidWlkIjogMSwNCiAgImRpc3BsYXlOYW1lIjogIllvdXIgTmFtZSIsDQogICJwaG90b1VSTCI6ICJodHRwOi8vZXhhbXBsZS5jb20iDQp9';
  const signature = 'TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ';
  return `${header}.${payload}.${signature}`;
};
