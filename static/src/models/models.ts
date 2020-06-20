export interface User {
  displayName: string;
  photoURL: string;
  uid: string;
}

export interface JWTHeader {
  typ: string;
  alg: string;
}

export interface JWTClaim {
  uid: string;
  displayName: string;
  photoURL: string;
}

export interface JWT {
  header: JWTHeader;
  claim: JWTClaim;
  signature: string;
}
