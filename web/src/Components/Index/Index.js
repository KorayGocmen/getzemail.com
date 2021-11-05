import React, {useState} from "react";
import { withRouter, useHistory } from "react-router-dom";

function Index(props) {
  const history = useHistory();
  const [address, setAddress] = useState(null);

  function handleSearch(){
    history.push(`/${address}`);
  }

  return (
    <div className="Index" style={{
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    }}>
      <form>
        <p>Search for any inbox</p>
        <input type="text" name="address" onChange={e =>setAddress(e.target.value)} />
        <button type="submit" onClick={handleSearch}> Search </button>
      </form>
    </div>
  );
}

export default withRouter(Index);
