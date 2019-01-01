using System.Collections;
using TMPro;
using UnityEngine;
using UnityEngine.Networking;
using UnityEngine.UI;

public class CubeController : MonoBehaviour
{
    public TMP_InputField Input;
    public Button RecordButton;
    public AudioSource Audio;

    public int RecordingTime = 3;

    public float Speed = 5;
    public float Angle = 0;
    public bool Moving = true;
    public bool Jumping = false;

    private void Start() {
        Input.onDeselect.AddListener(Deselect);
        RecordButton.onClick.AddListener(Record);
    }

    void Update () {
       if(Moving) {
            transform.position += transform.forward.normalized * Speed * Time.deltaTime;

            transform.Rotate(new Vector3(0, Angle, 0));
        }
    }

    private void Record() {
        StartCoroutine(_Record());
    }

    private IEnumerator _Record() {

        Audio.clip = Microphone.Start(null, true, RecordingTime, 44100);
        Audio.Play();

        yield return new WaitForSeconds(RecordingTime);

        Debug.Log(Audio.clip.length);

        byte[] audio = WavUtility.FromAudioClip(Audio.clip);

        string uri = "http://167.99.186.233/get_intent/Me";

        StartCoroutine(_SendPostRequest(uri, audio));
    }

    private void TriggerCommand(string command) {
        switch(command)
        {
            case "run":
                Moving = true;
                Speed = 10;
                break;
            case "walk":
                Moving = true;
                Speed = 5;
                break;
            case "jump":
                Jumping = true;
                break;
            case "stop":
                Moving = false;
                Jumping = false;
                Speed = 5;
                break;
            case "forwards":
                Speed = Mathf.Abs(Speed);
                break;
            case "backwards":
                Speed = -Mathf.Abs(Speed);
                break;
            case "left":
                Angle = -1;
                break;
            case "right":
                Angle = 1;
                break;
            
        }
    }

    private void Deselect(string text) {
        GetText(Input.text);
    }

    private IEnumerator _SendPostRequest(string uri, byte[] data) {
        using (var request = UnityWebRequest.Put(uri, data))
        {
            yield return request.SendWebRequest();

            if (request.isHttpError || request.isNetworkError)
            {
                Debug.Log(request.error);
            }
            else if (request.responseCode == 200)
            {
                string responseBody = request.downloadHandler.text;
                Debug.Log("Got:" + responseBody);

                var result = JsonUtility.FromJson<APIResult>(responseBody);

                TriggerCommand(result.result);
            }
        }
    }

    public void GetText(string text) {
        if(string.IsNullOrEmpty(text)) { return; }

        string uri = string.Format("http://167.99.186.233/get_intent/Me/{0}", UnityWebRequest.EscapeURL(text));

        StartCoroutine(_SendGetRequest(uri));
    }

    private IEnumerator _SendGetRequest(string uri) {
        using (var request = UnityWebRequest.Get(uri))
        {
            yield return request.SendWebRequest();

            if(request.isHttpError || request.isNetworkError)
            {
                Debug.Log(request.error);
            }
            else if (request.responseCode == 200)
            {
                string responseBody = request.downloadHandler.text;
                Debug.Log("Got:" + responseBody);

                var result = JsonUtility.FromJson<APIResult>(responseBody);

                TriggerCommand(result.result);
            }
        }
    }

    private class APIResult {
        public string result;
    }
}
