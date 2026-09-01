import {useParams} from "react-router-dom";
import ReactPlayer from "react-player/youtube";
import "./StreamMovie.css";

const StreamMovie = () => {
  let params = useParams();
  let key = params.yt_id; // Youtube ID from DB & URL
  
  return (
    <div className="react-player-container">
      {(key != null) ? (
        <ReactPlayer 
          controls="true" 
          playing={true} 
          url={`https://www.youtube.com/watch?v=${key}`}
          width="100%" height="100%" />
      ): (
        null
      )}
    </div>
  )
}

export default StreamMovie;